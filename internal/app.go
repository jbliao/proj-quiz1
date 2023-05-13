package internel

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jbliao/proj-quiz1/internal/model"
	"github.com/jbliao/proj-quiz1/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	listen       int
	dbConnString string
	commentSvc   service.CommentService
}

func NewApp(listenPort int, sqlConnectString string) *App {
	return &App{
		listen:       listenPort,
		dbConnString: sqlConnectString,
	}
}

func errWrapper(err interface{}) gin.H {
	return gin.H{"error": fmt.Sprint(err)}
}

// Start the app
func (app *App) Run() error {

	// Init db connection
	db, err := gorm.Open(mysql.Open(app.dbConnString), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&model.Comment{})

	// Init services
	app.commentSvc = service.NewPersistedCommentService(db)

	// Configure gin routes and starting to run!
	r := gin.Default()

	comment := r.Group("quiz/v1/comment")
	{
		comment.POST("", app.createComment)
		comment.GET("/:uuid", app.getCommentByUUID)
		comment.PUT("/:uuid", app.updateCommentByUUID)
		comment.DELETE("/:uuid", app.deleteCommentByUUID)
	}

	if err := r.Run(fmt.Sprintf("0.0.0.0:%v", app.listen)); err != nil {
		return err
	}

	return nil
}

func (app *App) createComment(c *gin.Context) {
	cmt := &model.Comment{}
	if err := c.ShouldBindJSON(cmt); err != nil {
		c.JSON(http.StatusBadRequest, errWrapper(err))
		return
	}

	if cmt.Uuid != "" {
		c.JSON(http.StatusBadRequest, errWrapper("request should not contain uuid"))
		return
	}

	// validate and normalize parentid
	if cmt.ParentID != "" {
		if _, err := uuid.Parse(cmt.ParentID); err != nil {
			c.JSON(http.StatusBadRequest, errWrapper("parentid not uuid"))
			return
		}
	}

	cmt.Uuid = uuid.New().String()
	if err := app.commentSvc.EnsureComment(cmt); err != nil {
		c.JSON(http.StatusInternalServerError, errWrapper(err))
		return
	}

	c.JSON(http.StatusOK, cmt)
}

func (app *App) getCommentByUUID(c *gin.Context) {
	givenUuid := c.Param("uuid")
	if _, err := uuid.Parse(givenUuid); err != nil {
		c.JSON(http.StatusBadRequest, errWrapper("cannot parse uuid"))
		return
	}

	comment, err := app.commentSvc.GetComment(givenUuid)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusBadRequest
		}
		c.JSON(status, errWrapper(err))
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (app *App) updateCommentByUUID(c *gin.Context) {
	givenUuid := c.Param("uuid")
	if _, err := uuid.Parse(givenUuid); err != nil {
		c.JSON(http.StatusBadRequest, errWrapper("cannot parse uuid"))
		return
	}

	newcmt := &model.Comment{}
	if err := c.ShouldBindJSON(newcmt); err != nil {
		c.JSON(http.StatusBadRequest, errWrapper(err))
		return
	}

	newcmt.Uuid = givenUuid
	if err := app.commentSvc.EnsureComment(newcmt); err != nil {
		c.JSON(http.StatusInternalServerError, errWrapper(err))
		return
	}

	c.JSON(http.StatusOK, newcmt)
	return
}

func (app *App) deleteCommentByUUID(c *gin.Context) {
	givenUuid := c.Param("uuid")
	if _, err := uuid.Parse(givenUuid); err != nil {
		c.JSON(http.StatusBadRequest, errWrapper("cannot parse uuid"))
		return
	}

	if err := app.commentSvc.DeleteComment(givenUuid); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusBadRequest
		}
		c.JSON(status, errWrapper(err))
		return
	}

	c.Status(http.StatusOK)
}

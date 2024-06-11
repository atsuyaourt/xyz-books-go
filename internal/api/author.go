package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Author struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
} //@name Author

func newAuthor(arg db.Author) Author {
	return Author{
		FirstName:  arg.FirstName,
		LastName:   arg.LastName,
		MiddleName: arg.LastName,
	}
}

type createAuthorReq struct {
	FirstName  string `json:"first_name" binding:"required,min=1"`
	LastName   string `json:"last_name" binding:"required,min=1"`
	MiddleName string `json:"middle_name" binding:"omitempty,min=1"`
} //@name CreateAuthorParams

// CreateAuthor
//
//	@Summary	Create author
//	@Tags		authors
//	@Accept		json
//	@Produce	json
//	@Param		req	body		createAuthorReq	true	"Create author parameters"
//	@Success	201	{object}	Author
//	@Router		/authors [post]
func (s *Server) CreateAuthor(ctx *gin.Context) {
	var req createAuthorReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAuthorParams(req)

	author, err := s.store.CreateAuthor(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, newAuthor(author))
}

type getAuthorReq struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

// GetAuthor
//
//	@Summary	Get author
//	@Tags		authors
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"author ID"
//	@Success	200	{object}	Author
//	@Router		/authors/{id} [get]
func (s *Server) GetAuthor(ctx *gin.Context) {
	var req getAuthorReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	author, err := s.store.GetAuthor(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("author not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newAuthor(author))
}

type listAuthorsReq struct {
	Page    int32 `form:"page,default=1" binding:"omitempty,min=1"`     // page number
	PerPage int32 `form:"per_page,default=5" binding:"omitempty,min=1"` // limit
} //@name ListAuthorsParams

type PaginatedAuthors = PaginatedList[Author] //@name PaginatedAuthors

// ListAuthors
//
//	@Summary	List authors
//	@Tags		authors
//	@Accept		json
//	@Produce	json
//	@Param		req	query		listAuthorsReq	false	"List authors parameters"
//	@Success	200	{object}	PaginatedAuthors
//	@Router		/authors [get]
func (s *Server) ListAuthors(ctx *gin.Context) {
	var req listAuthorsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	offset := (req.Page - 1) * req.PerPage

	arg := db.ListAuthorsParams{
		Limit:  int64(req.PerPage),
		Offset: int64(offset),
	}
	authors, err := s.store.ListAuthors(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	numAuthors := len(authors)
	items := make([]Author, numAuthors)
	for i, author := range authors {
		items[i] = newAuthor(author)
	}

	count, err := s.store.CountAuthors(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := NewPaginatedList[Author](req.Page, req.PerPage, int32(count), items)

	ctx.JSON(http.StatusOK, res)
}

type updateAuthorUri struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

type updateAuthorReq struct {
	FirstName  string `json:"first_name" binding:"omitempty,min=1"`
	LastName   string `json:"last_name" binding:"omitempty,min=1"`
	MiddleName string `json:"middle_name" binding:"omitempty,min=1"`
} //@name UpdateAuthorParams

// UpdateAuthor
//
//	@Summary	Update author
//	@Tags		authors
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int				true	"author ID"
//	@Param		req	body		updateAuthorReq	true	"Update author parameters"
//	@Success	200	{object}	Author
//	@Router		/authors/{id} [put]
func (s *Server) UpdateAuthor(ctx *gin.Context) {
	var uri updateAuthorUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateAuthorReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAuthorParams{
		AuthorID: uri.ID,
		FirstName: sql.NullString{
			String: req.FirstName,
			Valid:  len(req.FirstName) > 0,
		},
		LastName: sql.NullString{
			String: req.LastName,
			Valid:  len(req.LastName) > 0,
		},
		MiddleName: sql.NullString{
			String: req.MiddleName,
			Valid:  len(req.MiddleName) > 0,
		},
	}

	author, err := s.store.UpdateAuthor(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("author not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newAuthor(author))
}

type deleteAuthorUri struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

// DeleteAuthor
//
//	@Summary	Delete author
//	@Tags		authors
//	@Accept		json
//	@Produce	json
//	@Param		id	path	int	true	"author ID"
//	@Success	204
//	@Router		/authors/{id} [delete]
func (s *Server) DeleteAuthor(ctx *gin.Context) {
	var req deleteAuthorUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.store.DeleteAuthor(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

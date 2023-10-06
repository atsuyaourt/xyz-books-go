package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	db "github.com/emiliogozo/xyz-books/db/sqlc"
	"github.com/emiliogozo/xyz-books/internal/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Book struct {
	Title           string   `json:"title"`
	ISBN13          string   `json:"isbn13"`
	ISBN10          string   `json:"isbn10"`
	Price           float64  `json:"price"`
	PublicationYear int64    `json:"publication_year"`
	ImageUrl        string   `json:"image_url"`
	Edition         string   `json:"edition"`
	Authors         []string `json:"authors"`
	Publisher       string   `json:"publisher"`
} //@name Book

type newBookArg struct {
	Book      db.Book
	Publisher string
	Authors   []string
}

func newBook(arg newBookArg) Book {
	res := Book{
		Title:           arg.Book.Title,
		Price:           arg.Book.Price,
		PublicationYear: arg.Book.PublicationYear,
		Publisher:       arg.Publisher,
		Authors:         arg.Authors,
	}

	if arg.Book.Isbn13.Valid {
		res.ISBN13 = arg.Book.Isbn13.String
	}
	if arg.Book.Isbn10.Valid {
		res.ISBN10 = arg.Book.Isbn10.String
	}

	if arg.Book.ImageUrl.Valid {
		res.ImageUrl = arg.Book.ImageUrl.String
	}
	if arg.Book.Edition.Valid {
		res.Edition = arg.Book.Edition.String
	}

	return res
}

type createBookReq struct {
	Book struct {
		Title           string  `json:"title" binding:"required"`
		ISBN13          string  `json:"isbn13" binding:"omitempty,len=13,required_without=Isbn10,isbn13"`
		ISBN10          string  `json:"isbn10" binding:"omitempty,len=10,required_without=Isbn13,isbn10"`
		Price           float64 `json:"price" binding:"required,numeric"`
		PublicationYear int64   `json:"publication_year" binding:"required,numeric,min=1000"`
		ImageUrl        string  `json:"image_url" binding:"omitempty,url"`
		Edition         string  `json:"edition" binding:"omitempty"`
	} `json:"book"`
	Authors   []string `json:"authors" binding:"required,min=1"`
	Publisher string   `json:"publisher" binding:"required"`
} //@name CreateBookParams

// CreateBook
//
//	@Summary	Create book
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		req	body		createBookReq	true	"Create book parameters"
//	@Success	201	{object}	Book
//	@Router		/books [post]
func (s *Server) CreateBook(ctx *gin.Context) {
	var req createBookReq
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var authors []util.Name
	for i := range req.Authors {
		n := util.NewName(req.Authors[i])
		if n.Valid() {
			authors = append(authors, *n)
		}
	}

	publisher := cases.Title(language.English, cases.Compact).String(req.Publisher)

	arg := db.CreateBookTxParams{
		Book: db.CreateBookParams{
			Title: req.Book.Title,
			Isbn13: sql.NullString{
				String: req.Book.ISBN13,
				Valid:  len(req.Book.ISBN13) == 13,
			},
			Isbn10: sql.NullString{
				String: req.Book.ISBN10,
				Valid:  len(req.Book.ISBN10) == 10,
			},
			Price:           req.Book.Price,
			PublicationYear: req.Book.PublicationYear,
			ImageUrl: sql.NullString{
				String: req.Book.ImageUrl,
				Valid:  len(req.Book.ImageUrl) > 0,
			},
			Edition: sql.NullString{
				String: req.Book.Edition,
				Valid:  len(req.Book.Edition) > 0,
			},
		},
		Publisher: publisher,
		Authors:   authors,
	}

	book, err := s.store.CreateBookTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := newBook(newBookArg{
		Book:      book,
		Authors:   req.Authors,
		Publisher: publisher,
	})

	ctx.JSON(http.StatusCreated, res)
}

type getBookReq struct {
	ISBN13 string `uri:"isbn" binding:"required,isbn13"`
}

// GetBook
//
//	@Summary	Get book
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		isbn	path		string	true	"ISBN-13"
//	@Success	200		{object}	Book
//	@Router		/books/{isbn} [get]
func (s *Server) GetBook(ctx *gin.Context) {
	var req getBookReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	isbn := util.NewISBN(req.ISBN13)

	book, err := s.store.GetBookByISBN(ctx, db.GetBookByISBNParams{
		Isbn13: sql.NullString{
			String: isbn.ISBN13,
			Valid:  true,
		},
		Isbn10: sql.NullString{
			String: isbn.ISBN10,
			Valid:  true,
		},
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("book not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newBook(newBookArg{
		Book:      book.Book,
		Authors:   strings.Split(book.Authors, ","),
		Publisher: book.PublisherName,
	}))
}

type listBooksReq struct {
	Title              string  `form:"title" binding:"omitempty"`
	MinPrice           float32 `form:"min_price,default=-1.0" binding:"omitempty,numeric"`
	MaxPrice           float32 `form:"max_price,default=-1.0" binding:"omitempty,numeric"`
	MinPublicationYear int32   `form:"min_publication_year,default=-1" binding:"omitempty,numeric"`
	MaxPublicationYear int32   `form:"max_publication_year,default=-1" binding:"omitempty,numeric"`
	Author             string  `form:"author" binding:"omitempty"`
	Publisher          string  `form:"publisher" binding:"omitempty"`
	Page               int32   `form:"page,default=1" binding:"omitempty,min=1"`            // page number
	PerPage            int32   `form:"per_page,default=5" binding:"omitempty,min=1,max=30"` // limit
} //@name ListBooksParams

type PaginatedBooks = PaginatedList[Book] //@name PaginatedBooks

// ListBooks
//
//	@Summary	List books
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		req	query		listBooksReq	false	"List books parameters"
//	@Success	200	{object}	PaginatedBooks
//	@Router		/books [get]
func (s *Server) ListBooks(ctx *gin.Context) {
	var req listBooksReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	offset := (req.Page - 1) * req.PerPage

	arg := db.ListBooksParams{
		Limit:  int64(req.PerPage),
		Offset: int64(offset),
		Title: sql.NullString{
			String: req.Title,
			Valid:  len(req.Title) > 0,
		},
		Author: sql.NullString{
			String: req.Author,
			Valid:  len(req.Author) > 0,
		},
		Publisher: sql.NullString{
			String: req.Publisher,
			Valid:  len(req.Publisher) > 0,
		},
		MinPrice: sql.NullFloat64{
			Float64: float64(req.MinPrice),
			Valid:   req.MinPrice >= 0,
		},
		MaxPrice: sql.NullFloat64{
			Float64: float64(req.MaxPrice),
			Valid:   req.MaxPrice > req.MinPrice,
		},
		MinPublicationYear: sql.NullInt64{
			Int64: int64(req.MinPublicationYear),
			Valid: req.MinPublicationYear > 999,
		},
		MaxPublicationYear: sql.NullInt64{
			Int64: int64(req.MaxPublicationYear),
			Valid: req.MaxPublicationYear > req.MinPublicationYear,
		},
	}
	bookRows, err := s.store.ListBooks(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	numBooks := len(bookRows)
	items := make([]Book, numBooks)
	for i, book := range bookRows {
		items[i] = newBook(newBookArg{
			Book:      book.Book,
			Authors:   strings.Split(book.Authors, ","),
			Publisher: book.PublisherName,
		})
	}

	count, err := s.store.CountBooks(ctx, db.CountBooksParams{
		Title:              arg.Title,
		Author:             arg.Author,
		Publisher:          arg.Publisher,
		MinPrice:           arg.MinPrice,
		MaxPrice:           arg.MaxPrice,
		MinPublicationYear: arg.MinPublicationYear,
		MaxPublicationYear: arg.MinPublicationYear,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := NewPaginatedList[Book](req.Page, req.PerPage, int32(count), items)

	ctx.JSON(http.StatusOK, res)
}

type updateBookUri struct {
	ISBN13 string `uri:"isbn" binding:"required,isbn13"`
}

type updateBookReq struct {
	Title           string  `json:"title" binding:"omitempty,min=1"`
	NewISBN13       string  `json:"isbn13" binding:"omitempty,isbn13"`
	NewISBN10       string  `json:"isbn10" binding:"omitempty,isbn10"`
	Price           float32 `json:"price" binding:"omitempty,numeric"`
	PublicationYear int32   `json:"publication_year"  binding:"omitempty,numeric"`
	ImageUrl        string  `json:"image_url"  binding:"omitempty,url"`
} //@name UpdateBookParams

// UpdateBook
//
//	@Summary	Update book
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		isbn	path		string			true	"ISBN-13"
//	@Param		req		body		updateBookReq	true	"Update book parameters"
//	@Success	200		{object}	Book
//	@Router		/books/{isbn} [put]
func (s *Server) UpdateBook(ctx *gin.Context) {
	var uri updateBookUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateBookReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	isbn := util.NewISBN(uri.ISBN13)

	arg := db.UpdateBookByISBNParams{
		Isbn13: sql.NullString{
			String: uri.ISBN13,
			Valid:  true,
		},
		Isbn10: sql.NullString{
			String: isbn.ISBN10,
			Valid:  true,
		},
		Title: sql.NullString{
			String: req.Title,
			Valid:  len(req.Title) > 0,
		},
		Price: sql.NullFloat64{
			Float64: float64(req.Price),
			Valid:   req.Price > 0,
		},
		PublicationYear: sql.NullInt64{
			Int64: int64(req.PublicationYear),
			Valid: req.PublicationYear > 999,
		},
		ImageUrl: sql.NullString{
			String: req.ImageUrl,
			Valid:  len(req.ImageUrl) > 0,
		},
	}

	if (len(req.NewISBN13) == 13) && (len(req.NewISBN10) == 10) {
		isbn1 := util.NewISBN(req.NewISBN13)
		isbn2 := util.NewISBN(req.NewISBN10)

		if isbn1.ISBN10 == isbn2.ISBN10 {
			arg.NewIsbn13 = sql.NullString{
				String: req.NewISBN13,
				Valid:  true,
			}
			arg.NewIsbn10 = sql.NullString{
				String: req.NewISBN10,
				Valid:  true,
			}
		}
	} else if len(req.NewISBN13) == 13 {
		arg.NewIsbn13 = sql.NullString{
			String: req.NewISBN13,
			Valid:  true,
		}
	} else if len(req.NewISBN10) == 10 {
		arg.NewIsbn10 = sql.NullString{
			String: req.NewISBN10,
			Valid:  true,
		}
	}

	_, err := s.store.UpdateBookByISBN(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("book not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	book, err := s.store.GetBookByISBN(ctx, db.GetBookByISBNParams{
		Isbn13: arg.Isbn13,
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("book not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newBook(newBookArg{
		Book:      book.Book,
		Authors:   strings.Split(book.Authors, ","),
		Publisher: book.PublisherName,
	}))
}

type deleteBookUri struct {
	ISBN13 string `uri:"isbn" binding:"required,isbn13"`
}

// DeleteBook
//
//	@Summary	Delete book
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		isbn	path	string	true	"ISBN-13"
//	@Success	204
//	@Router		/books/{isbn} [delete]
func (s *Server) DeleteBook(ctx *gin.Context) {
	var req deleteBookUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.store.DeleteBookByISBN(ctx, db.DeleteBookByISBNParams{
		Isbn13: sql.NullString{
			String: req.ISBN13,
			Valid:  true,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

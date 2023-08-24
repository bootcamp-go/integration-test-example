package handlers

import (
	"app/internal/sellers"
	"app/internal/sellers/repository"
	"app/internal/sellers/storage"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

// Integration test for ControllerSellers.GetById with Default Implementation for RepositorySellers
func TestIntegration_ControllerSellers_GetByID(t *testing.T) {
	t.Run("found seller by id", func(t *testing.T) {
		// arrange
		st := storage.NewStorageSellersMock()
		st.ReadFunc = func() (s map[int]*sellers.Seller, err error) {
			s = map[int]*sellers.Seller{
				1: {
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "johndoe@gmail.com",
				},
			}
			return
		}

		rp := repository.NewRepositorySellersDefault(st)
		impl := NewControllerSellers(rp)
		hd := impl.GetById()

		// act
		inputR := httptest.NewRequest("GET", "/sellers/1", nil)
		
		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", "1")

		inputR = inputR.WithContext(
			context.WithValue(
				inputR.Context(),
				chi.RouteCtxKey,
				rCtx,
			),
		)
		// inputR.ctx = context.WithValue(
		// 	inputR.ctx,
		// 	"id",
		// 	"1",
		// )
		inputW := httptest.NewRecorder()
		hd(inputW, inputR)

		// assert
		outputCode := http.StatusOK
		outputBody := `{"message":"success","data":{"id":1,"first_name":"John","last_name":"Doe","email":"johndoe@gmail.com"},"error":false}`
		require.Equal(t, outputCode, inputW.Code)
		require.JSONEq(t, outputBody, inputW.Body.String())
		require.Equal(t, 1, st.Calls.Read)
		require.Equal(t, 0, st.Calls.Write)
	})

	t.Run("not found seller by id", func(t *testing.T) {
		// arrange
		st := storage.NewStorageSellersMock()
		st.ReadFunc = func() (s map[int]*sellers.Seller, err error) {
			s = make(map[int]*sellers.Seller)
			return
		}

		rp := repository.NewRepositorySellersDefault(st)
		impl := NewControllerSellers(rp)
		hd := impl.GetById()

		// act
		inputR := httptest.NewRequest("GET", "/sellers/1", nil)

		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", "1")

		inputR = inputR.WithContext(
			context.WithValue(
				inputR.Context(),
				chi.RouteCtxKey,
				rCtx,
			),
		)
		inputW := httptest.NewRecorder()
		hd(inputW, inputR)

		// assert
		outputCode := http.StatusNotFound
		outputBody := `{"message":"seller not found","data":null,"error":true}`
		require.Equal(t, outputCode, inputW.Code)
		require.JSONEq(t, outputBody, inputW.Body.String())
		require.Equal(t, 1, st.Calls.Read)
		require.Equal(t, 0, st.Calls.Write)
	})

	t.Run("storage error", func(t *testing.T) {
		// arrange
		st := storage.NewStorageSellersMock()
		st.ReadFunc = func() (s map[int]*sellers.Seller, err error) {
			err = storage.ErrStorageSellersInternal
			return
		}

		rp := repository.NewRepositorySellersDefault(st)
		impl := NewControllerSellers(rp)
		hd := impl.GetById()

		// act
		inputR := httptest.NewRequest("GET", "/sellers/1", nil)

		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", "1")

		inputR = inputR.WithContext(
			context.WithValue(
				inputR.Context(),
				chi.RouteCtxKey,
				rCtx,
			),
		)
		inputW := httptest.NewRecorder()
		hd(inputW, inputR)

		// assert
		outputCode := http.StatusInternalServerError
		outputBody := `{"message":"internal error","data":null,"error":true}`
		require.Equal(t, outputCode, inputW.Code)
		require.JSONEq(t, outputBody, inputW.Body.String())
		require.Equal(t, 1, st.Calls.Read)
		require.Equal(t, 0, st.Calls.Write)
	})

	t.Run("bad request - path param not presented", func(t *testing.T) {
		// arrange
		st := storage.NewStorageSellersMock()
		
		rp := repository.NewRepositorySellersDefault(st)

		impl := NewControllerSellers(rp)
		hd := impl.GetById()

		// act
		inputR := httptest.NewRequest("GET", "/sellers/", nil)
		inputW := httptest.NewRecorder()
		hd(inputW, inputR)

		// assert
		outputCode := http.StatusBadRequest
		outputBody := `{"message":"bad request","data":null,"error":true}`
		require.Equal(t, outputCode, inputW.Code)
		require.JSONEq(t, outputBody, inputW.Body.String())
		require.Equal(t, 0, st.Calls.Read)
		require.Equal(t, 0, st.Calls.Write)
	})

	t.Run("bad request - path param not integer", func(t *testing.T) {
		// arrange
		st := storage.NewStorageSellersMock()

		rp := repository.NewRepositorySellersDefault(st)

		impl := NewControllerSellers(rp)
		hd := impl.GetById()

		// act
		inputR := httptest.NewRequest("GET", "/sellers/abc", nil)

		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", "abc")

		inputR = inputR.WithContext(
			context.WithValue(
				inputR.Context(),
				chi.RouteCtxKey,
				rCtx,
			),
		)
		inputW := httptest.NewRecorder()
		hd(inputW, inputR)

		// assert
		outputCode := http.StatusBadRequest
		outputBody := `{"message":"bad request","data":null,"error":true}`
		require.Equal(t, outputCode, inputW.Code)
		require.JSONEq(t, outputBody, inputW.Body.String())
		require.Equal(t, 0, st.Calls.Read)
		require.Equal(t, 0, st.Calls.Write)
	})
}

func TestIntegration_ControllerSellers_Save(t *testing.T) {
	t.Run("success to save seller", func(t *testing.T) {
		// arrange
		st := storage.NewStorageSellersMock()
		st.ReadFunc = func() (s map[int]*sellers.Seller, err error) {
			s = make(map[int]*sellers.Seller)
			return
		}
		st.WriteFunc = func(s map[int]*sellers.Seller) (err error) {
			return
		}

		rp := repository.NewRepositorySellersDefault(st)
		impl := NewControllerSellers(rp)
		hd := impl.Save()

		// act
		inputR := httptest.NewRequest("POST", "/sellers/", io.NopCloser(
			strings.NewReader(
				`{"first_name":"John","last_name":"Doe","email":"johndoe@gmail.com"}`,
			),
		))
		inputW := httptest.NewRecorder()
		hd(inputW, inputR)

		// assert
		outputCode := http.StatusCreated
		outputBody := `{"message":"success","data":{"id":1,"first_name":"John","last_name":"Doe","email":"johndoe@gmail.com"}, "error":false}`
		require.Equal(t, outputCode, inputW.Code)
		require.JSONEq(t, outputBody, inputW.Body.String())
		require.Equal(t, 1, st.Calls.Read)
		require.Equal(t, 1, st.Calls.Write)
	})

	t.Run("bad request - empty body", func(t *testing.T) {
		// arrange
		st := storage.NewStorageSellersMock()

		rp := repository.NewRepositorySellersDefault(st)
		impl := NewControllerSellers(rp)
		hd := impl.Save()

		// act
		inputR := httptest.NewRequest("POST", "/sellers/", nil)
		inputW := httptest.NewRecorder()
		hd(inputW, inputR)

		// assert
		outputCode := http.StatusBadRequest
		outputBody := `{"message":"bad request","data":null,"error":true}`
		require.Equal(t, outputCode, inputW.Code)
		require.JSONEq(t, outputBody, inputW.Body.String())
		require.Equal(t, 0, st.Calls.Read)
		require.Equal(t, 0, st.Calls.Write)
	})

	t.Run("storage error - write failed", func(t *testing.T) {
		// arrange
		st := storage.NewStorageSellersMock()
		st.ReadFunc = func() (s map[int]*sellers.Seller, err error) {
			s = make(map[int]*sellers.Seller)
			return
		}
		st.WriteFunc = func(s map[int]*sellers.Seller) (err error) {
			err = storage.ErrStorageSellersInternal
			return
		}

		rp := repository.NewRepositorySellersDefault(st)
		impl := NewControllerSellers(rp)
		hd := impl.Save()

		// act
		inputR := httptest.NewRequest("POST", "/sellers/", io.NopCloser(
			strings.NewReader(
				`{"first_name":"John","last_name":"Doe","email":"johndoe@gmail.com"}`,
			),
		))
		inputW := httptest.NewRecorder()
		hd(inputW, inputR)

		// assert
		outputCode := http.StatusInternalServerError
		outputBody := `{"message":"internal error","data":null,"error":true}`
		require.Equal(t, outputCode, inputW.Code)
		require.JSONEq(t, outputBody, inputW.Body.String())
		require.Equal(t, 1, st.Calls.Read)
		require.Equal(t, 1, st.Calls.Write)
	})
}
package gateway

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type unavailableEngine struct{}

func NewUnavailableEngine() Engine { return unavailableEngine{} }
func (unavailableEngine) Complete(_ *gin.Context, _ string, _ Request) (Completion, error) {
	return Completion{}, errors.New("provider unavailable")
}

package helper

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ListQuery struct {
	Page   int
	Limit  int
	Offset int
	Search string
}

func ParamID(c *fiber.Ctx, key string) (int, error) {
	val, err := strconv.Atoi(c.Params(key))
	if err != nil || val <= 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "invalid parameter "+key)
	}
	return val, nil
}

func ParseListQuery(c *fiber.Ctx) ListQuery {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return ListQuery{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
		Search: c.Query("search", ""),
	}
}
package route

import (
	deckUseCase "api/application/deck"
	"api/application/detail"
	"api/application/search"
	searchDeckUseCase "api/application/search/deck"
	meiliQueryService "api/infrastructure/meilisearch/query_service"
	mysqlQueryService "api/infrastructure/mysql/query_service"
	"api/infrastructure/mysql/repository"
	deckPre "api/presentation/deck"
	detailPre "api/presentation/detail"
	searchPre "api/presentation/search"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339_nano}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	v1 := e.Group("/v1")

	cardSearchRoute(v1)
	cardDetailRoute(v1)
	deckRoute(v1)
}

func cardSearchRoute(g *echo.Group) {
	// Meilisearchの実装に切り替え
	pokemonRepository := meiliQueryService.NewPokemonQueryService()
	trainerRepository := meiliQueryService.NewTrainerQueryService()
	energyRepository := meiliQueryService.NewEnergyQueryService()
	searchRepository := search.NewSearchPokemonAndTrainerUseCase(
		pokemonRepository,
		trainerRepository,
		energyRepository,
	)
	deckQueryService := meiliQueryService.NewDeckQueryService()
	searchDeckUseCase := searchDeckUseCase.NewSearchDeckUseCase(deckQueryService)
	h := searchPre.NewSearchHandler(searchRepository, searchDeckUseCase)

	group := g.Group("/search")
	group.GET("/cards", h.SearchCardList)
	group.GET("/decks", h.SearchDeckList)
}

func cardDetailRoute(g *echo.Group) {
	detailRepository := mysqlQueryService.NewDetailQueryService()
	detailUseCase := detail.NewFetchDetailUseCase(detailRepository)
	h := detailPre.NewDetailHandler(detailUseCase)

	group := g.Group("/cards")
	group.GET("/detail/:card_type/:id", h.FetchDetail)
}

func deckRoute(g *echo.Group) {
	deckRepository := repository.NewDeckRepository()
	cardRepository := repository.NewCardRepository()

	listDeckUseCase := deckUseCase.NewListDeckUseCase(deckRepository)
	createDeckUseCase := deckUseCase.NewCreateDeckUseCase(deckRepository, cardRepository)
	validateDeckUseCase := deckUseCase.NewValidateDeckUseCase(cardRepository)
	updateDeckUseCase := deckUseCase.NewUpdateDeckUseCase(deckRepository, cardRepository)
	deleteDeckUseCase := deckUseCase.NewDeleteDeckUseCase(deckRepository)

	deckHandler := deckPre.NewDeckHandler(
		listDeckUseCase,
		createDeckUseCase,
		validateDeckUseCase,
		updateDeckUseCase,
		deleteDeckUseCase,
	)

	group := g.Group("/decks")
	group.GET("", deckHandler.GetAllDecks)
	group.GET("/detail/:id", deckHandler.GetDeckById)
	group.POST("/create", deckHandler.CreateDeck)
	group.POST("/validate", deckHandler.ValidateDeck)
	group.POST("/edit/:id", deckHandler.UpdateDeck)
	group.DELETE("/delete/:id", deckHandler.DeleteDeck)
}

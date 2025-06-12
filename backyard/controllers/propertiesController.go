package controllers

import (
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
    "github.com/hospedate/backyard/services"
)

type PropertiesController interface {
    Create(NewPropertyRequest *models.NewPropertyRequest) (*models.Property, error)
    GetById(id string) *models.Property
    SearchProperties(queryParams *models.PropertiesSearchParams) []*models.Property
    CreateBlock(newBlockRequest *models.NewBlockRequest) (*models.Block, error)
    GetBlocksByPropertyId(id_property string) []*models.Block
    DeleteBlockById(id_property string, id_block string) (int64, error)
    Edit(propertyEditParams models.PropertyEditRequest) (*models.Property, error)
}

type PropertiesControllerImpl struct {
    repository         repositories.PropertiesRepository
    ordersRepo         repositories.OrdersRepository
    airbnbFetcher      services.AirbnbFetcher
    usersCreditService services.UsersCreditService
}

func NewPropertiesController(
    repository repositories.PropertiesRepository,
    ordersRepo repositories.OrdersRepository,
    airbnbFetcher services.AirbnbFetcher,
    usersCreditService services.UsersCreditService,
) PropertiesController {

    return &PropertiesControllerImpl{
        repository:         repository,
        ordersRepo:         ordersRepo,
        airbnbFetcher:      airbnbFetcher,
        usersCreditService: usersCreditService,
    }
}

func (c *PropertiesControllerImpl) Create(newPropertyRequest *models.NewPropertyRequest) (*models.Property, error) {
    property, err := c.repository.Save(newPropertyRequest)

    if err == repositories.ErrDuplicateKey {
        return nil, ErrDuplicateKey
    } else if err == repositories.ErrUserDoesNotExist {
        return nil, ErrUserDoesNotExist
    } else if err != nil {
        return nil, UnknownError
    }

    // trigger airbnb fetcher
    c.airbnbFetcher.SendPropertyEvent(property)

    return property, nil
}

func (c *PropertiesControllerImpl) GetById(id string) *models.Property {
    return c.repository.GetById(id)
}

func (c *PropertiesControllerImpl) SearchProperties(searchParams *models.PropertiesSearchParams) []*models.Property {
    return c.repository.Search(searchParams)
}

func (c *PropertiesControllerImpl) CreateBlock(newBlockRequest *models.NewBlockRequest) (*models.Block, error) {
    block, err := c.repository.SaveBlock(newBlockRequest)

    if err == repositories.ErrCollision {
        return nil, ErrCollision
    } else if err != nil {
        return nil, UnknownError
    }

    return block, nil
}

func (c *PropertiesControllerImpl) GetBlocksByPropertyId(id_property string) []*models.Block {
    return c.repository.GetBlocksByPropertyId(id_property)
}

func (c *PropertiesControllerImpl) DeleteBlockById(id_property string, id_block string) (int64, error) {
    return c.repository.DeleteBlockById(id_property, id_block)
}

func (c *PropertiesControllerImpl) Edit(propertyEditParams models.PropertyEditRequest) (*models.Property, error) {
    property := c.repository.GetById(propertyEditParams.Id)

    if property == nil {
        return nil, ErrPropertyDoesNotExist
    }

    if property.Status == "archived" {
        return nil, ErrPropertyArchivedStatus
    }

    if shouldArchive(propertyEditParams) {
        return c.archiveProperty(property, propertyEditParams)
    }

    propertyEdited, err := c.repository.Edit(propertyEditParams)
    if err != nil {
        return nil, err
    }

    if !property.IsVerified && propertyEdited.IsVerified {
        c.usersCreditService.HandleCreditOwnerEvent(property.UserId)
    }

    return propertyEdited, nil
}

func shouldArchive(propertyEditParams models.PropertyEditRequest) bool {
    return propertyEditParams.Status != nil && *propertyEditParams.Status == "archived"
}

func (c *PropertiesControllerImpl) archiveProperty(property *models.Property, propertyEditParams models.PropertyEditRequest) (*models.Property, error) {
    if err := c.failIfActiveOrders(property); err != nil {
        return nil, err
    }

    return c.repository.SavePropertyArchived(property, &propertyEditParams)
}

func (c *PropertiesControllerImpl) failIfActiveOrders(property *models.Property) error {
    invalidStatuses := []string{"pending", "confirmed", "ephemeral", "in_progress"}

    orderSearch := models.OrdersSearchParams{OwnerId: &property.UserId}
    orders := c.ordersRepo.Search(&orderSearch)

    for _, order := range orders {
        if order.PropertyId == property.Id {
            for _, invalidStatus := range invalidStatuses {
                if order.Status == invalidStatus {
                    return ErrPropertyHasActiveOrders
                }
            }
        }
    }

    return nil
}

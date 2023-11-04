package usecase

import (
	"fmt"

	"github.com/walnuts1018/openchokin/back/domain"
)

type MoneyProviderSummary struct {
	ID      string
	Name    string
	Balance float64
}
type MoneyProvidersSummaryResponse struct {
	Providers []MoneyProviderSummary
}

func (u Usecase) GetMoneyProvidersSummary(userID string) (MoneyProvidersSummaryResponse, error) {
	moneyProviders, err := u.db.GetMoneyProvidersByUserID(userID)
	if err != nil {
		return MoneyProvidersSummaryResponse{}, err
	}

	var providersSummary []MoneyProviderSummary
	for _, provider := range moneyProviders {
		balance, err := u.db.GetMoneyPoolBalance(provider.ID, false)
		if err != nil {
			// handle the error according to your error policy
			// for example, you could log it and continue with the next provider
			continue
		}
		providersSummary = append(providersSummary, MoneyProviderSummary{
			ID:      provider.ID,
			Name:    provider.Name,
			Balance: balance,
		})
	}

	return MoneyProvidersSummaryResponse{Providers: providersSummary}, nil
}

type MoneyProviderResponse struct {
	ID        string
	Name      string
	CreatorID string
	Balance   float64
}

// UpdateMoneyProvider updates an existing money provider's name and balance.
func (u Usecase) UpdateMoneyProvider(userID string, moneyProviderID string, name string, balance float64) (MoneyProviderResponse, error) {
	// Get the existing MoneyProvider to check if it belongs to the user.
	existingProvider, err := u.db.GetMoneyProvider(moneyProviderID)
	if err != nil {
		return MoneyProviderResponse{}, err
	}

	// Check if the user is authorized to update the MoneyProvider.
	if existingProvider.CreatorID != userID {
		return MoneyProviderResponse{}, fmt.Errorf("unauthorized to update money provider: %s", moneyProviderID)
	}

	// Update the MoneyProvider details.
	updatedProvider := domain.MoneyProvider{
		ID:        moneyProviderID,
		Name:      name,
		CreatorID: userID,
		Balance:   balance,
	}

	// Use the DB interface method to update.
	err = u.db.UpdateMoneyProvider(updatedProvider)
	if err != nil {
		return MoneyProviderResponse{}, err
	}

	// Return the updated MoneyProvider response.
	return MoneyProviderResponse{
		ID:        updatedProvider.ID,
		Name:      updatedProvider.Name,
		CreatorID: updatedProvider.CreatorID,
		Balance:   updatedProvider.Balance,
	}, nil
}

// AddMoneyProvider adds a new money provider for a user with the given name and balance.
func (u Usecase) AddMoneyProvider(userID string, name string, balance float64) (MoneyProviderResponse, error) {
	// Create a new MoneyProvider instance.
	newProvider := domain.MoneyProvider{
		Name:      name,
		CreatorID: userID,
		Balance:   balance,
	}

	// Use the DB interface method to create a new MoneyProvider.
	createdProvider, err := u.db.NewMoneyProvider(newProvider)
	if err != nil {
		return MoneyProviderResponse{}, err
	}

	// Return the response with the new MoneyProvider details.
	return MoneyProviderResponse{
		ID:        createdProvider.ID,
		Name:      createdProvider.Name,
		CreatorID: createdProvider.CreatorID,
		Balance:   createdProvider.Balance,
	}, nil
}

// DeleteMoneyProvider deletes an existing money provider.
func (u Usecase) DeleteMoneyProvider(userID string, moneyProviderID string) error {
	// Get the existing MoneyProvider to check if it belongs to the user.
	provider, err := u.db.GetMoneyProvider(moneyProviderID)
	if err != nil {
		return err
	}

	// Check if the user is authorized to delete the MoneyProvider.
	if provider.CreatorID != userID {
		return fmt.Errorf("unauthorized to delete money provider: %s", moneyProviderID)
	}

	// Use the DB interface method to delete the MoneyProvider.
	err = u.db.DeleteMoneyProvider(moneyProviderID)
	if err != nil {
		return err
	}

	return nil
}

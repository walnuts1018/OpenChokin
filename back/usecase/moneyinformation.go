package usecase

import (
	"time"

	"github.com/walnuts1018/openchokin/back/domain"
)

type MoneySumResponse struct {
	MoneyProviderSum       float64
	ActualMoneyPoolSum     float64
	ForecastedMoneyPoolSum float64
}

func (u Usecase) GetMoneyInformation(userID string, loginUserID string) (MoneySumResponse, error) {
	var response MoneySumResponse

	// Retrieve all MoneyPools associated with the user.
	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		return response, err
	}

	// Process each MoneyPool based on access rights.
	for _, pool := range moneyPools {
		var hasAccess bool // Flag to check access permissions

		// Direct access for the same user or public access for others.
		if userID == loginUserID {
			hasAccess = true // Users have full access to their own pools
		} else if loginUserID != "" {
			// Check shared access or public type for logged-in users.
			shared, err := u.db.IsMoneyPoolSharedWithUser(pool.ID, loginUserID)
			if err != nil {
				return response, err // Error checking shared status
			}
			hasAccess = shared || pool.Type == domain.PublicTypePublic
		} else {
			// No login user ID provided; only include public pools.
			hasAccess = pool.Type == domain.PublicTypePublic
		}

		// If the user has access, sum up the actual and forecasted balances.
		if hasAccess {
			balance, err := u.db.GetMoneyPoolBalance(pool.ID, false)
			if err != nil {
				return response, err
			}
			response.ActualMoneyPoolSum += balance

			forecastedBalance, err := u.db.GetMoneyPoolBalance(pool.ID, true)
			if err != nil {
				return response, err
			}
			response.ForecastedMoneyPoolSum += forecastedBalance
		}
	}

	// Retrieve all MoneyProviders for the user and calculate the sum.
	moneyProviders, err := u.db.GetMoneyProvidersByUserID(userID)
	if err != nil {
		return response, err
	}
	for _, provider := range moneyProviders {
		response.MoneyProviderSum += provider.Balance
	}

	return response, nil
}

func (u Usecase) GetMoneyInformationOfDate(userID string, loginUserID string, date time.Time) (MoneySumResponse, error) {
	var response MoneySumResponse

	// Retrieve all MoneyProviders for the user and calculate the sum of their balances.
	moneyProviders, err := u.db.GetMoneyProvidersByUserID(userID)
	if err != nil {
		return response, err
	}
	for _, provider := range moneyProviders {
		response.MoneyProviderSum += provider.Balance
	}

	// Retrieve MoneyPools based on access rights
	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		return response, err
	}
	for _, pool := range moneyPools {
		// Calculate balances only if the user has access to the money pool
		var hasAccess bool
		if userID == loginUserID {
			hasAccess = true // User has full access to their own pools
		} else if loginUserID != "" {
			// Check if the money pool is shared or public
			shared, err := u.db.IsMoneyPoolSharedWithUser(pool.ID, loginUserID)
			if err != nil {
				return response, err // Error checking shared status
			}
			hasAccess = shared || pool.Type == domain.PublicTypePublic
		} else {
			// No login user ID provided; public pools only
			hasAccess = pool.Type == domain.PublicTypePublic
		}

		if hasAccess {
			// Calculate actual balance up to the given date
			balance, err := u.db.GetMoneyPoolBalanceOfDate(pool.ID, date, false)
			if err != nil {
				return response, err
			}
			response.ActualMoneyPoolSum += balance

			// Calculate forecasted balance up to the given date
			balance, err = u.db.GetMoneyPoolBalanceOfDate(pool.ID, date, true)
			if err != nil {
				return response, err
			}
			response.ForecastedMoneyPoolSum += balance
		}
	}

	return response, nil
}

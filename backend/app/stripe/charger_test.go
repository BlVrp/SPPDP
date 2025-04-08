package stripe_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"one-help/app/stripe"
	"one-help/internal/logger/zaplog"
)

func TestCharger(t *testing.T) {
	t.Skip("for manual use")

	charger := stripe.NewCharger(zaplog.NewLog(), stripe.Config{
		SecretAPIKey:   "",
		RedirectDomain: "",
		PriceID:        "",
	})

	t.Run("SetupChargeSessionURL", func(t *testing.T) {
		url, id, err := charger.SetupChargeSessionURL("/fundraises/donation/id")
		require.NoError(t, err)
		fmt.Println(id)
		fmt.Println(url)
	})

	t.Run("GetSessionPaymentData", func(t *testing.T) {
		paid, funded, err := charger.GetSessionPaymentData("cs_test_a1Qt7394rR4TdfzCyBpx1La116znKYoWlCVIBMV5tfRlUnXsWymknP4fTn")
		require.NoError(t, err)
		fmt.Println(paid)
		fmt.Println(funded)
	})
}

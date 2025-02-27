package securityhub

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	// Maximum amount of time to wait for an AdminAccount to return Enabled
	adminAccountEnabledTimeout = 5 * time.Minute

	// Maximum amount of time to wait for an AdminAccount to return NotFound
	adminAccountNotFoundTimeout = 5 * time.Minute

	standardsSubscriptionCreateTimeout = 3 * time.Minute
	standardsSubscriptionDeleteTimeout = 3 * time.Minute
)

// waitAdminAccountEnabled waits for an AdminAccount to return Enabled
func waitAdminAccountEnabled(ctx context.Context, conn *securityhub.SecurityHub, adminAccountID string) (*securityhub.AdminAccount, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{adminStatusNotFound},
		Target:  []string{securityhub.AdminStatusEnabled},
		Refresh: statusAdminAccountAdmin(ctx, conn, adminAccountID),
		Timeout: adminAccountEnabledTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*securityhub.AdminAccount); ok {
		return output, err
	}

	return nil, err
}

// waitAdminAccountNotFound waits for an AdminAccount to return NotFound
func waitAdminAccountNotFound(ctx context.Context, conn *securityhub.SecurityHub, adminAccountID string) (*securityhub.AdminAccount, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{securityhub.AdminStatusDisableInProgress},
		Target:  []string{adminStatusNotFound},
		Refresh: statusAdminAccountAdmin(ctx, conn, adminAccountID),
		Timeout: adminAccountNotFoundTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*securityhub.AdminAccount); ok {
		return output, err
	}

	return nil, err
}

func waitStandardsSubscriptionCreated(ctx context.Context, conn *securityhub.SecurityHub, arn string) (*securityhub.StandardsSubscription, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{securityhub.StandardsStatusPending},
		Target:  []string{securityhub.StandardsStatusReady, securityhub.StandardsStatusIncomplete},
		Refresh: statusStandardsSubscription(ctx, conn, arn),
		Timeout: standardsSubscriptionCreateTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*securityhub.StandardsSubscription); ok {
		return output, err
	}

	return nil, err
}

func waitStandardsSubscriptionDeleted(ctx context.Context, conn *securityhub.SecurityHub, arn string) (*securityhub.StandardsSubscription, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{securityhub.StandardsStatusDeleting},
		Target:  []string{standardsStatusNotFound, securityhub.StandardsStatusIncomplete},
		Refresh: statusStandardsSubscription(ctx, conn, arn),
		Timeout: standardsSubscriptionDeleteTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*securityhub.StandardsSubscription); ok {
		return output, err
	}

	return nil, err
}

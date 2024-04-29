package hcloud

import (
	"context"
	"fmt"
	"time"
)

// WatchOverallProgress watches several actions' progress until they complete
// with success or error. This watching happens in a goroutine and updates are
// provided through the two returned channels:
//
//   - The first channel receives percentage updates of the progress, based on
//     the number of completed versus total watched actions. The return value
//     is an int between 0 and 100.
//   - The second channel returned receives errors for actions that did not
//     complete successfully, as well as any errors that happened while
//     querying the API.
//
// By default, the method keeps watching until all actions have finished
// processing. If you want to be able to cancel the method or configure a
// timeout, use the [context.Context]. Once the method has stopped watching,
// both returned channels are closed.
//
// WatchOverallProgress uses the [WithPollBackoffFunc] of the [Client] to wait
// until sending the next request.
func (c *ActionClient) WatchOverallProgress(ctx context.Context, actions []*Action) (<-chan int, <-chan error) {
	errCh := make(chan error, len(actions))
	progressCh := make(chan int)

	go func() {
		defer close(errCh)
		defer close(progressCh)

		completedIDs := make([]int64, 0, len(actions))
		watchIDs := make(map[int64]struct{}, len(actions))
		for _, action := range actions {
			watchIDs[action.ID] = struct{}{}
		}

		retries := 0
		previousProgress := 0

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-time.After(c.action.client.pollBackoffFunc(retries)):
				retries++
			}

			opts := ActionListOpts{}
			for watchID := range watchIDs {
				opts.ID = append(opts.ID, watchID)
			}

			as, err := c.AllWithOpts(ctx, opts)
			if err != nil {
				errCh <- err
				return
			}
			if len(as) == 0 {
				// No actions returned for the provided IDs, they do not exist in the API.
				// We need to catch and fail early for this, otherwise the loop will continue
				// indefinitely.
				errCh <- fmt.Errorf("failed to wait for actions: remaining actions (%v) are not returned from API", opts.ID)
				return
			}

			progress := 0
			for _, a := range as {
				switch a.Status {
				case ActionStatusRunning:
					progress += a.Progress
				case ActionStatusSuccess:
					delete(watchIDs, a.ID)
					completedIDs = append(completedIDs, a.ID)
				case ActionStatusError:
					delete(watchIDs, a.ID)
					completedIDs = append(completedIDs, a.ID)
					errCh <- fmt.Errorf("action %d failed: %w", a.ID, a.Error())
				}
			}

			progress += len(completedIDs) * 100
			if progress != 0 && progress != previousProgress {
				sendProgress(progressCh, progress/len(actions))
				previousProgress = progress
			}

			if len(watchIDs) == 0 {
				return
			}
		}
	}()

	return progressCh, errCh
}

// WatchProgress watches one action's progress until it completes with success
// or error. This watching happens in a goroutine and updates are provided
// through the two returned channels:
//
//   - The first channel receives percentage updates of the progress, based on
//     the progress percentage indicated by the API. The return value is an int
//     between 0 and 100.
//   - The second channel receives any errors that happened while querying the
//     API, as well as the error of the action if it did not complete
//     successfully, or nil if it did.
//
// By default, the method keeps watching until the action has finished
// processing. If you want to be able to cancel the method or configure a
// timeout, use the [context.Context]. Once the method has stopped watching,
// both returned channels are closed.
//
// WatchProgress uses the [WithPollBackoffFunc] of the [Client] to wait until
// sending the next request.
func (c *ActionClient) WatchProgress(ctx context.Context, action *Action) (<-chan int, <-chan error) {
	errCh := make(chan error, 1)
	progressCh := make(chan int)

	go func() {
		defer close(errCh)
		defer close(progressCh)

		retries := 0

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-time.After(c.action.client.pollBackoffFunc(retries)):
				retries++
			}

			a, _, err := c.GetByID(ctx, action.ID)
			if err != nil {
				errCh <- err
				return
			}
			if a == nil {
				errCh <- fmt.Errorf("failed to wait for action %d: action not returned from API", action.ID)
				return
			}

			switch a.Status {
			case ActionStatusRunning:
				sendProgress(progressCh, a.Progress)
			case ActionStatusSuccess:
				sendProgress(progressCh, 100)
				errCh <- nil
				return
			case ActionStatusError:
				errCh <- a.Error()
				return
			}
		}
	}()

	return progressCh, errCh
}

// sendProgress allows the user to only read from the error channel and ignore any progress updates.
func sendProgress(progressCh chan int, p int) {
	select {
	case progressCh <- p:
		break
	default:
		break
	}
}

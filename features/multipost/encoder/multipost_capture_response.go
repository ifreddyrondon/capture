package encoder

import "github.com/ifreddyrondon/capture/features/multipost/decoder"

const (
	replyForEmail            = "we will send an email when finished"
	replyForCallback         = "we will notify via callback url (POST) when finished"
	replyForEmailAndCallback = "we will send an email and notify via callback url (POST) when finished"
	defaultReply             = "you will see it when finished"
)

type MultiPOSTCaptureResponse struct {
	IgnoreErrors      bool   `json:"ignore_errors"`
	CapturesToProcess int    `json:"captures_to_process"`
	NotificationEmail string `json:"notification_email"`
	CallbackURL       string `json:"callback_url"`
	Message           string `json:"message"`
}

func NewMultiPOSTCaptureResponse(captures decoder.MultiPOSTCaptures) MultiPOSTCaptureResponse {
	res := MultiPOSTCaptureResponse{
		IgnoreErrors:      captures.IgnoreErrors,
		CapturesToProcess: len(captures.Captures),
		NotificationEmail: captures.Notifications.Email,
		CallbackURL:       captures.Notifications.CallbackURL,
		Message:           defaultReply,
	}

	if captures.Notifications.Email != "" && captures.Notifications.CallbackURL == "" {
		res.Message = replyForEmail
	} else if captures.Notifications.CallbackURL != "" && captures.Notifications.Email == "" {
		res.Message = replyForCallback
	} else if captures.Notifications.CallbackURL != "" && captures.Notifications.Email != "" {
		res.Message = replyForEmailAndCallback
	}

	return res
}

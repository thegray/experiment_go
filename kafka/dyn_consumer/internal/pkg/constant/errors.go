package constant

const (
	ErrorCodeUnknown            int = 480
	ErrorCodeNoOrder            int = 481
	ErrorCodeInvalidCompanyCode int = 484
	ErrorCodeNoVA               int = 485
	ErrorCodeInvalidConnection  int = 489
	ErrorCodeInvalidServiceReq  int = 497

	ErrorCodeBNICancelAlreadyInvalid  string = "103"
	ErrorCodeBNISystemOffline         string = "997"
	ErrorCodeBNISystemUnexpectedError string = "009"
	ErrorCodeBNIErrCode3Times         int    = 1000

	ErrorCodeInvalidPrefix int = 499
	ErrorCodeCommon        int = 402

	ErrorCodeCommonInvalid    int = 1001
	ErrorCodeWrongAmountOrder int = 483
	ErrorCodeExpiredOrder     int = 482
	ErrorCodePaidOrder        int = 488
	ErrorCodeCancelledOrder   int = 486
	ErrorCodeRefundedOrder    int = 487
)

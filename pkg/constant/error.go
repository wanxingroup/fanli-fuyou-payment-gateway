package constant

const ErrorCodeInternalError = 500000

const ErrorCodeFuYouMerchantIdRequired = 413001 // 富友商户号为必填
const ErrorMessageFuYouMerchantIdRequired = "富友商户号为必填"

const ErrorCodeAppIdRequired = 413002 // AppId 必填
const ErrorMessageAppIdRequired = "AppId 必填"

const ErrorCodeOpenIdRequired = 413003 // OpenId 必填
const ErrorMessageOpenIdRequired = "OpenId 必填"

const ErrorCodePayTypeRequired = 413004 // 支付类型为必填
const ErrorMessagePayTypeRequired = "支付类型为必填"

const ErrorCodePayTypeNotSupportOption = 413005 // 不支持输入的支付类型选项
const ErrorMessagePayTypeNotSupportOption = "不支持输入的支付类型选项"

const ErrorCodeDescriptionRequired = 413006 // 商品描述为必填
const ErrorMessageDescriptionRequired = "商品描述为必填"

const ErrorCodeDescriptionTooLarge = 413007 // 商品描述为必填
const ErrorMessageDescriptionTooLarge = "商品描述过长"

const ErrorCodeMerchantIdRequired = 413008 // 商家 ID 为必填
const ErrorMessageMerchantIdRequired = "商家 ID 为必填"

const ErrorCodeTransactionNotExist = 413009 // 业务流水不存在
const ErrorMessageTransactionNotExist = "业务流水不存在"

const ErrorCodeCloseTransactionResponseError = 513001 // 请求富友关闭支付流水时，富友返回失败，详情查看错误信息
const ErrorMessageCloseTransactionResponseError = "请求富友关闭支付流水时，富友返回失败，详情查看错误信息"

const ErrorCodeConvertErrorCodeFailed = 513002 // 转换富友错误号时，转换失败（可能不是整形）
const ErrorMessageConvertErrorCodeFailed = "转换富友错误号时，转换失败（可能不是整形）"

# Swagger Golang REST API demo
This is a sample API server

## Version: 1.0.0


### /accounts

#### POST
##### Summary:

Add a new account

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body | Account object that needs to be added to the store | Yes | [AccountRequest](#accountrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [AccountResponce](#accountresponce) |
| 400 | invalid input: cannot decode |  |

#### GET
##### Summary:

Get all accounts

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [AccountResponce](#accountresponce) ] |
| 500 | internal server error |  |

### /accounts/{accountId}/payments

#### GET
##### Summary:

Get account payments

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| accountId | path | ID of account | Yes | int64 |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [AccountPaymentResponse](#accountpaymentresponse) ] |
| 500 | internal server error |  |

### /payments

#### POST
##### Summary:

Do a new payment

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body | Payment object that needs to be added to the store | Yes | [DoPaymentRquest](#dopaymentrquest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [AccountResponce](#accountresponce) |
| 400 | invalid input: cannot decode |  |

### Models


#### AccountRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string | the new account name | Yes |
| currency | string | the account currency | Yes |
| balance | int64 | the account initial balance | Yes |

#### AccountResponce

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | int64 | the account ID | Yes |
| name | string |  | Yes |
| currency | string |  | Yes |
| balance | int64 |  | Yes |

#### DoPaymentRquest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| amount | int64 | how much money need to transfer | Yes |
| to_id | string | transfer to account ID | Yes |
| from_id | string | transfer from account ID | Yes |

#### AccountPaymentResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | int64 | the payment ID | Yes |
| amount | int64 | how much money transferred | Yes |
| to_account | string | transferred to account ID | Yes |
| from_account | string | transferred from account ID | Yes |
| direction | string | transfer direction | Yes |

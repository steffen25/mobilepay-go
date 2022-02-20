package mobilepay

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"strconv"
	"testing"
)

var config = &Config{
	HTTPClient: newDefaultHTTPClient(),
	URL:        TestBaseUrl,
}

func TestPayments_List(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	testdata, err := ioutil.ReadFile("testdata/list_payments.json")
	if err != nil {
		t.Fatal(err)
	}

	testdata = bytes.Replace(testdata, []byte("PAGE_SIZE"), []byte(strconv.Itoa(10)), 1)
	testdata = bytes.Replace(testdata, []byte("NEXT_PAGE_NUMBER"), []byte(strconv.Itoa(2)), 1)

	gock.New(TestBaseUrl).
		Get("/v1/payments").
		MatchParam("pageNumber", "1").
		MatchParam("pageSize", "10").
		Reply(200).
		JSON(testdata)

	client := New("test", "test", config)
	ctx := context.TODO()

	listOptions := ListOptions{
		PageSize:   10,
		PageNumber: 1,
	}

	data, err := client.Payment.Get(ctx, listOptions)

	assert.Nil(t, err)
	assert.True(t, len(data.Payments) == 3)
	assert.Exactly(t, 10, data.PageSize)
	assert.Exactly(t, 2, data.NextPageNumber)
}

func TestPayments_Get(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	testdata, err := ioutil.ReadFile("testdata/get_payment.json")
	if err != nil {
		t.Fatal(err)
	}

	testdata = bytes.Replace(testdata, []byte("PAYMENT_ID"), []byte("186d2b31-ff25-4414-9fd1-bfe9807fa8b7"), 1)

	gock.New(TestBaseUrl).
		Get("/v1/payments/186d2b31-ff25-4414-9fd1-bfe9807fa8b7").
		Reply(200).
		JSON(testdata)

	client := New("test", "test", config)
	ctx := context.TODO()

	payment, err := client.Payment.Find(ctx, "186d2b31-ff25-4414-9fd1-bfe9807fa8b7")
	assert.Nil(t, err)
	assert.NotNil(t, payment.PaymentId)
}

func TestPayments_Get_Empty_PaymentId(t *testing.T) {
	client := New("test", "test", config)
	ctx := context.TODO()

	payment, err := client.Payment.Find(ctx, "")
	assert.Error(t, err)
	assert.IsType(t, &ArgError{}, err)
	assert.Nil(t, payment)
}

func TestPayments_Create(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	testdata, err := ioutil.ReadFile("testdata/create_payment.json")
	if err != nil {
		t.Fatal(err)
	}

	testdata = bytes.Replace(testdata, []byte("PAYMENT_ID"), []byte("186d2b31-ff25-4414-9fd1-bfe9807fa8b7"), 2)

	gock.New(TestBaseUrl).
		Post("/v1/payments").
		Reply(200).
		JSON(testdata)

	client := New("test", "test", config)
	ctx := context.TODO()

	params := &PaymentParams{
		Amount:         1050,
		IdempotencyKey: "7347ba06-95c5-4181-82e5-7c7a23609a0e",
		PaymentPointId: "1f8ed17f-f310-4f40-a7a4-df78185efbdd",
		RedirectUri:    "app://callback",
		Reference:      "test",
		Description:    "this is a test payment",
	}

	payment, err := client.Payment.Create(ctx, params)

	assert.Nil(t, err)
	assert.Equal(t, "186d2b31-ff25-4414-9fd1-bfe9807fa8b7", payment.PaymentId)
	assert.Equal(t, "mobilepay://merchant_payments?payment_id=186d2b31-ff25-4414-9fd1-bfe9807fa8b7", payment.MobilePayAppRedirectUri)
}

//
func TestPayments_Cancel(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New(TestBaseUrl).
		Post("/v1/payments/186d2b31-ff25-4414-9fd1-bfe9807fa8b7/cancel").
		Reply(204)

	client := New("test", "test", config)
	ctx := context.TODO()

	err := client.Payment.Cancel(ctx, "186d2b31-ff25-4414-9fd1-bfe9807fa8b7")
	assert.Nil(t, err)
}

//
func TestPayments_Capture(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New(TestBaseUrl).
		Post("/v1/payments/206d2b31-ff25-4414-9fd1-bfe9807fa8b7/capture").
		Reply(204)

	client := New("test", "test", config)
	ctx := context.TODO()

	err := client.Payment.Capture(ctx, "206d2b31-ff25-4414-9fd1-bfe9807fa8b7", 100)
	assert.Nil(t, err)
}

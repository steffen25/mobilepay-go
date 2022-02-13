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

func TestRefund_List(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	testdata, err := ioutil.ReadFile("testdata/list_refunds.json")
	if err != nil {
		t.Fatal(err)
	}

	testdata = bytes.Replace(testdata, []byte("PAGE_SIZE"), []byte(strconv.Itoa(10)), 1)
	testdata = bytes.Replace(testdata, []byte("NEXT_PAGE_NUMBER"), []byte(strconv.Itoa(2)), 1)

	gock.New(TestBaseUrl).
		Get("/v1/refunds").
		MatchParam("pageNumber", "1").
		MatchParam("pageSize", "10").
		Reply(200).
		JSON(testdata)

	client := NewClient("test", "test", config)

	ctx := context.TODO()

	opts := &RefundsListOptions{
		ListOptions: ListOptions{PageSize: 10, PageNumber: 1},
		PaymentId:   "211444eb-1c4e-4194-a58f-905d97877cc5",
	}

	data, err := client.Payment.Refund.List(ctx, opts)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(data.Refunds))
	assert.Equal(t, 10, data.PageSize)
	assert.Equal(t, 2, data.NextPageNumber)

}

func TestRefund_Create(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	testdata, err := ioutil.ReadFile("testdata/create_refund.json")
	if err != nil {
		t.Fatal(err)
	}

	testdata = bytes.Replace(testdata, []byte("PAYMENT_ID"), []byte("211444eb-1c4e-4194-a58f-905d97877cc5"), 1)
	testdata = bytes.Replace(testdata, []byte("DESCRIPTION"), []byte("this is a test payment"), 1)
	testdata = bytes.Replace(testdata, []byte("REFERENCE"), []byte("Test"), 1)
	testdata = bytes.Replace(testdata, []byte("AMOUNT"), []byte(strconv.Itoa(100)), 1)

	gock.New(TestBaseUrl).
		Post("/v1/refunds").
		Reply(200).
		JSON(testdata)

	client := NewClient("test", "test", config)
	ctx := context.TODO()

	params := &RefundParams{
		Amount:         100,
		IdempotencyKey: "7576910d-9789-4fef-a72e-877d89afec94",
		PaymentId:      "211444eb-1c4e-4194-a58f-905d97877cc5",
		Reference:      "test",
		Description:    "this is a test payment",
	}

	refund, err := client.Payment.Refund.Create(ctx, params)

	assert.Nil(t, err)
	assert.Equal(t, "211444eb-1c4e-4194-a58f-905d97877cc5", refund.PaymentId)
	assert.Equal(t, "this is a test payment", refund.Description)
	assert.Equal(t, "Test", refund.Reference)
	assert.Equal(t, 100, refund.Amount)
}

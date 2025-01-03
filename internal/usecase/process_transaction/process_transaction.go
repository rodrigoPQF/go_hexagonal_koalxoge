package process_transaction

import (
	"github.com/rodrigoPQF/go_hexagonal_koalxoge/internal/entity"
)

type ProcessTransaction struct {
	Repository entity.TransactionRepository
}

func NewProcessTransaction(repository entity.TransactionRepository) *ProcessTransaction {
	return &ProcessTransaction{Repository: repository}
}

func (p *ProcessTransaction) Execute(input TransactionDtoInput) (TransactionDtoOutput, error) {
	transaction := entity.NewTransaction()
	transaction.ID = input.ID
	transaction.AccountID = input.AccountID
	transaction.Amount = input.Amount
	invalidTransaction := transaction.IsValid()

	if invalidTransaction != nil {
		return p.rejectedTransaction(transaction, invalidTransaction)
	}
	return p.approveTransaction(transaction)
}

func (p *ProcessTransaction) approveTransaction(transaction *entity.Transaction) (TransactionDtoOutput, error) {
	err := p.Repository.Insert(transaction.ID, transaction.AccountID, transaction.Amount, "approved", "")
	if err != nil {
		return TransactionDtoOutput{}, err
	}
	output2 := TransactionDtoOutput{
		ID: transaction.ID,
		Status: "approved",
		ErrorMessage: "",
	}

	return output2, nil
}

func (p *ProcessTransaction) rejectedTransaction(transaction *entity.Transaction, invalidTransaction error) (TransactionDtoOutput, error) {
	err := p.Repository.Insert(transaction.ID, transaction.AccountID, transaction.Amount, "rejected", invalidTransaction.Error())
	if err != nil {
		return TransactionDtoOutput{}, err
	}
	output := TransactionDtoOutput{
		ID: transaction.ID,
		Status: "rejected",
		ErrorMessage: invalidTransaction.Error(),
	}
	return output, nil
}
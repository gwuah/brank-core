module Swipe
    module Business
      class UpdateMonoAccount < ApplicationInteractor
        parameters :account_id
  
        def call
          update_account_balance
          update_latest_transactions
        end
  
        private
  
        def update_latest_transactions
          transactions = Mono.transactions(mono_account.mono_id)
          return unless transactions
  
          mono_account.update(transactions: transactions[0..9])
        end
  
        def update_account_balance
          result = Mono.find_account(mono_account.mono_id)
          return unless result
          
          mono_account.update(balance: result["balance"].to_d, balance_updated_at: DateTime.now)
        end
  
        def mono_account
          MonoAccount.find(account_id)
        end
      end
    end
  end
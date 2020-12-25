module Swipe
    module Data
      class LinkBankAccount < ApplicationInteractor
        parameters :business, :code
  
        def call
          account_id = Mono.enroll(code: code)
          fail!(:enrollment_code, "Your enrollment code is invalid") unless account_id
  
          account = Mono.find_account(account_id)
          fail!(:account, "Failed to link bank account") unless account
          
          mono_account = business.mono_accounts.find_by_account_number(account["accountNumber"])
  
          if mono_account.nil?
            mono_account = business.mono_accounts.
              create!(account_number: account["accountNumber"], account_name:
              account["name"], currency: account["currency"], bank_name:
              account["institution"]["name"], bank_code:
              account["institution"]["code"], account_type: account["type"],
              balance: account["balance"], mono_id: account["_id"])
          end
  
          context.mono_account = mono_account
        end
      end
    end
  end
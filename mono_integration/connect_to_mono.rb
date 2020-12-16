module Mutations
    class ConnectToMono < Mutations::BaseMutation
      argument :code, String, required: true
  
      field :mono_account, Types::MonoAccount, null: true
  
      def execute(code:)
        ensure_authorized!
  
        result = Swipe::Data::LinkBankAccount.call!(business: @current_user.business, code: code)
        respond 200, mono_account: result.mono_account
      end
    end
  end
class UpdateMonoAccountBalanceWorker
    include Sidekiq::Worker
  
    def perform
      MonoAccount.where("balance_updated_at::date = ? OR balance_updated_at IS NULL", 1.day.ago.to_date).find_each do |mono_account|
        Swipe::Business::UpdateMonoAccount.call!(account_id: mono_account.id)
      end
    end
  end
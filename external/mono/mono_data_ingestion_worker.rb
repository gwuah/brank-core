class MonoDataIngestionWorker
    include Sidekiq::Worker
  
    def perform
      MonoAccount.find_each do |account|
        AsyncInteractorWorker.perform_async("Swipe::Data::IngestBankStatements", account_id: account.id)
      end
    end
  end
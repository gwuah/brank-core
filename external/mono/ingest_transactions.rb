module Swipe
    module Data
      class IngestTransactions < ApplicationInteractor
        parameters :transaction
  
        def call
          return unless ENV["APP_ENV"] == "production"
          return unless transaction.is_a? BusinessPayment
          return unless transaction.completed?
  
          context.params = params = {
            sourceID: source_id,
            businessID: business_id,
            amount: amount,
            currencyISO: currency,
            type: type,
            narration: narration,
            created_at: created_at,
          }
  
          response = client.post("/chomp", params.to_json)
          fail!(:ingestion, "Failed to ingest") unless response.status == 200
        end
  
        delegate :amount, :created_at, to: :transaction
  
        private
  
        def source_id
          transaction.id
        end
  
        def business_id
          transaction.business_id
        end
  
        def narration
          transaction.business_payment_type.name.downcase.to_sym
        end
  
        def type
          transaction.direction.downcase
        end
  
        def currency
          transaction.business.currency.code
        end
  
        def client
          @@client ||= Faraday.new(url: ENV["RECON_BASE_URL"],
            headers: { Authorization: ENV["RECON_SECRET_KEY"], "Content-Type": "application/json" })
        end
      end
    end
  end
class Mono
    CLIENT = Faraday.new(url: "https://api.withmono.com",
      headers: { "Content-Type": "application/json",
      "mono-sec-key": ENV["MONO_PRIVATE_KEY"] })
  
    class << self
      def enroll(code:)
        response = CLIENT.post("/account/auth", { code: code }.to_json)
        return nil if response.status != 200
  
        JSON.parse(response.body)["id"]
      end
  
      def find_account(id)
        response = CLIENT.get("/accounts/#{id}")
        return nil if response.status != 200
  
        JSON.parse(response.body)["account"]
      end
  
      def bank_statements(account_id:, from_id:)
        response = CLIENT.get("/accounts/#{account_id}/statement?from_id=#{from_id.to_s}")
        return nil if response.status != 200
  
        JSON.parse(response.body)["data"]
      end
  
      def transactions(account_id)
        response = CLIENT.get("/accounts/#{account_id}/transactions")
        return nil if response.status != 200
  
        JSON.parse(response.body)["data"]
      end
    end
  end
  
# ObmennikValut_Grps-Rest
1. Запустить .yaml в deployments (поднятие бд)
2. Запустить main файл в gw-currency-wallet
3. Запустить main файл в gw-exchanger
4. существующие роуты и их примеры:

   POST: http://localhost:8080/api/v1/wallet/exchange
   {
   "from_currency": "USD",
   "to_currency": "RUB",
   "amount": 100,
   "user_id": 1
   }

   GET: http://localhost:8080/api/v1/balance

   POST: http://localhost:8080/api/v1/register
   {
   "username": "Ivan",
   "email": "Ivan@example.com",
   "password": "password"
   }

   POST: http://localhost:8080/api/v1/login
   {
   "currency": "usd",
   "amount": 100.0
   }

   POST: http://localhost:8080/api/v1/wallet/withdraw
   {
   "currency": "usd",
   "amount": 50.0
   }


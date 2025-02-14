TOKEN1=$(curl -s -X POST http://localhost:8080/auth \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser1", "password": "testpassword"}' | jq -r '.token')

#echo "TOKEN: $TOKEN"
#
TOKEN2=$(curl -s -X POST http://localhost:8080/auth \
     -H "Content-Type: application/json" \
     -d '{"username": "andrey", "password": "testpassword"}' | jq -r '.token')

#
TOKEN3=$(curl -s -X POST http://localhost:8080/auth \
     -H "Content-Type: application/json" \
     -d '{"username": "petya", "password": "testpassword"}' | jq -r '.token')

curl -s -X GET http://localhost:8080/info \
     -H "Authorization: Bearer $TOKEN2" \
     -H "Content-Type: application/json"

echo "\n"

curl -s -X POST http://localhost:8080/sendCoin \
     -H "Authorization: Bearer $TOKEN1" \
     -H "Content-Type: application/json" \
     -d '{"toUser": "andrey", "amount": 100}'

echo "\n"


curl -s -X GET http://localhost:8080/buy/pen \
     -H "Authorization: Bearer $TOKEN2" \
     -H "Content-Type: application/json"

echo "\n"

curl -s -X GET http://localhost:8080/info \
    -H "Authorization: Bearer $TOKEN2" \
    -H "Content-Type: application/json"

#echo "\n"
#
#curl -s -X GET http://localhost:8080/info \
#    -H "Authorization: Bearer $TOKEN1" \
#    -H "Content-Type: application/json"
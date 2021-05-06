Trabalho Prático 1 - Sistemas Distríbuidos 
Pedro Fratini Chem

#### Executando ####

Dentro da pasta raiz execute:
- go run app.go 127.0.0.1:5001 127.0.0.1:7001 <ip:porta> <ip:porta> ...
- Exemplo: go run app.go 127.0.0.1:5001 127.0.0.1:6001 127.0.0.1:7001  

- Importante! 
Os ips/porta: 127.0.0.1:5001 e 127.0.0.1:7001 devem obrigatoriamente fazer parte do conjunto de argumentos, pois eles 
vão representar o processo de delay. 


#### Validação ####

- Após a execução do programa serão gerados arquivos textos com o output de cada processo.
- Para validar cada output navague para a pasta "validator" e execute o comando:
- go run validador.go ../<nome_arquvivo>
- Exemplo: go run validador.go ../5001_output.txt


#### Erros ####
- Caso o programa gere algum erro inesperado basta executa-lo novamente.
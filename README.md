# go-expert-stress-test

Uma CLI simples para executar requests HTTP GET para uma URL

Baseado no `Apache Benchmark` que utiliza as letras `ab` em seu comando, o componente deste repo se chama `tb` the `Thiago Benchmark` :tada:

O repositório conta com um `Dockerfile` pronto para fornecer uma imagem e permitir o uso da cli diretamente via docker.

Para construir a imagem, execute o seguinte comando:

```bash
docker build -t tb .
```

Para executar a CLI utilizando docker, execute o seguinte comando:

```bash
docker run tb --url "http://example.com"
docker run tb --url "http://google.com"
```

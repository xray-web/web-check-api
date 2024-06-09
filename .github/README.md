<h1 align="center">Web Check API</h1>
<p align="center">
  <a href="https://github.com/lissy93/web-check">
    <img width="72" src="./web-check.png?raw=true" />
    <br />
  </a>
  <i>A light-weight Go API for discovering website data</i><br />
  <b><a href="https://web-check.xyz">Web Check</a> - <i>Gives you Xray Vision for any Website</i></b>
</p>

> [!NOTE]
> This is a very early work in progress, and is not yet feature complete or production ready.
> Stay tuned!

---

## Usage

### Developing

#### Getting Started

You will need [git](https://git-scm.com/) and [go](https://go.dev/) installed.
Then clone the repo and download dependencies.

```
git clone git@github.com:xray-web/web-check-api.git
cd web-check-api
go mod download
```

#### Start Server

```
make run
```

#### Run Tests

```
make test
```


### Deploying

#### Option 1: From Source

Follow the setup instructions above. Then build the binaries.
Then execute the output executable directly (e.g. `./bin/app`)

```
make build
```

#### Option 2: From Docker

```
docker run -p 8080:8080 lissy93/web-check-api
```

#### Option 3: Download Executable
From the releases tab, download the compiled binary for your system, and execute it.

---

## License

> _**[Web Check](https://github.com/Lissy93/web-check)** is licensed under [MIT](https://github.com/xray-web/web-check-api/blob/HEAD/LICENSE) © [Alicia Sykes](https://aliciasykes.com) 2024._<br>
> <sup align="right">For information, see <a href="https://tldrlegal.com/license/mit-license">TLDR Legal > MIT</a></sup>

<details>
<summary>Expand License</summary>

```
The MIT License (MIT)
Copyright (c) Alicia Sykes <alicia@omg.com> 

Permission is hereby granted, free of charge, to any person obtaining a copy 
of this software and associated documentation files (the "Software"), to deal 
in the Software without restriction, including without limitation the rights 
to use, copy, modify, merge, publish, distribute, sub-license, and/or sell 
copies of the Software, and to permit persons to whom the Software is furnished 
to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included install 
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANT ABILITY, FITNESS FOR A
PARTICULAR PURPOSE AND NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```

</details>


<!-- License + Copyright -->
<p  align="center">
  <i>© <a href="https://aliciasykes.com">Alicia Sykes</a> 2024</i><br>
  <i>Licensed under <a href="https://gist.github.com/Lissy93/143d2ee01ccc5c052a17">MIT</a></i><br>
  <a href="https://github.com/lissy93"><img src="https://i.ibb.co/4KtpYxb/octocat-clean-mini.png" /></a><br>
  <sup>Thanks for visiting :)</sup>
</p>

<!-- Dinosaurs are Awesome -->
<!-- 
                        . - ~ ~ ~ - .
      ..     _      .-~               ~-.
     //|     \ `..~                      `.
    || |      }  }              /       \  \
(\   \\ \~^..'                 |         }  \
 \`.-~  o      /       }       |        /    \
 (__          |       /        |       /      `.
  `- - ~ ~ -._|      /_ - ~ ~ ^|      /- _      `.
              |     /          |     /     ~-.     ~- _
              |_____|          |_____|         ~ - . _ _~_-_
-->



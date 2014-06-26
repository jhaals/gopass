package main

import (
  "net/http"

  "code.google.com/p/go-uuid/uuid"
  "github.com/bradfitz/gomemcache/memcache"

  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"

  "github.com/jhaals/gopass/crypt"
  "github.com/jhaals/gopass/random"
)


func main() {

  m := martini.Classic()
  mc := memcache.New("127.0.0.1:11211")
  m.Use(render.Renderer(render.Options{
    Layout: "layout",
  }))

  // index
  m.Get("/", func(r render.Render) {
    r.HTML(200, "index", nil)
  })

  m.Get("/404", func(r render.Render) {
    r.HTML(404, "404", nil)
  })

  // show secret
  m.Get("/:uuid/:key", func(params martini.Params, r render.Render) string {
    it, err := mc.Get(params["uuid"])

    if err != nil {
      r.Redirect("/404")
    }

    secret := crypt.Decrypt([]byte(params["key"]), string(it.Value))
    mc.Delete(params["uuid"])
    return secret
  })

  // No decryption key submitted in URL, ask user for it
  m.Get("/:uuid", func(params martini.Params, r render.Render, req *http.Request) {
    url := req.URL.Query()
    decryption_key := url.Get("p")
    // decryption key of valid length submitted.
    if len(decryption_key) >= 16 {
      r.Redirect("/" + params["uuid"] + "/" + decryption_key)
    }

    r.HTML(200, "get_secret", nil)
  })

  // Save secret
  m.Post("/", func(p *http.Request,r render.Render) {

    valid_lifetime := map[string]int32 {
        "1h": 3600,
        "1d": 3600 * 24,
        "1w": 3600 * 24 * 7,
    }
    secret := p.FormValue("secret")
    lifetime := p.FormValue("lifetime")

    if len(secret) >= 1000 {
        r.HTML(413, "failure", "This site is meant to store secrets not novels")
        return
    }

    if valid_lifetime[lifetime] == 0 {
      r.HTML(400, "failure", "Not a valid lifetime")
      return
    }

    uuid := uuid.New()
    decryption_key := random.RandomString(16)
    encrypted_secret := crypt.Encrypt([]byte(decryption_key), secret)
    err := mc.Set(&memcache.Item{
      Key: uuid,
      Value: []byte(encrypted_secret),
      Expiration: valid_lifetime[lifetime],
    })

    if err != nil {
      r.HTML(500, "failure", "Failed to save secret")
      return
    }

    url := "http://127.0.0.1:3000/" + uuid + "/" + decryption_key
    r.HTML(200, "secret_url", url)
  })

  m.Run()
}

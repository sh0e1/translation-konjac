package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
)

// NewAuthConfigure ...
func NewAuthConfigure(id, secret string) *AuthConfigure {
	return &AuthConfigure{
		channelID:     id,
		channelSecret: secret,
	}
}

// AuthConfigure ...
type AuthConfigure struct {
	channelID     string
	channelSecret string
}

// Auth ...
func Auth(config *AuthConfigure) middlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			item, err := memcache.Get(ctx, channelTokenCacheKey)
			if err != nil && err != memcache.ErrCacheMiss {
				log.Errorf(ctx, "%+v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err == memcache.ErrCacheMiss {
				u := url.Values{}
				u.Add("grant_type", "client_credentials")
				u.Add("client_id", config.channelID)
				u.Add("client_secret", config.channelSecret)

				client := urlfetch.Client(ctx)
				resp, err := client.PostForm("https://api.line.me/v2/oauth/accessToken", u)
				if err != nil || resp.StatusCode != http.StatusOK {
					log.Errorf(ctx, "%+v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer resp.Body.Close()

				var channelToken channelToken
				if err := json.NewDecoder(resp.Body).Decode(&channelToken); err != nil {
					log.Errorf(ctx, "%+v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				item = &memcache.Item{
					Key:        channelTokenCacheKey,
					Value:      []byte(channelToken.AccessToken),
					Expiration: time.Duration(channelToken.ExpiresIn-(60*60*24)) * time.Second,
				}
				if err := memcache.Set(ctx, item); err != nil {
					log.Errorf(ctx, "%v", err)
				}
			}

			token := string(item.Value)
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, channelTokenContextKey, token)))
		}
	}
}

// GetChannelTokenFromContext ...
func GetChannelTokenFromContext(ctx context.Context) string {
	token, _ := ctx.Value(channelTokenContextKey).(string)
	return token
}

const channelTokenCacheKey = "channel_token"

type channelToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Schema needed to validate our config in some kinda sane way without
// the need of a parser. https://jsonschema.net does a nice job of giving
// you something you can work with quickly. You can also read more about it
// on https://json-schema.org/understanding-json-schema/index.html where they
// do a pretty good job at showing you how everything works.
//
// A lot of the enums will need to have "" also added to them unless they are
// required in the config. This is because we are taking YAML packing it into
// a go struct and converting it to JSON.
package main

var schema = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "http://example.com/root.json",
  "type": "object",
  "title": "The Root Schema",
  "required": [
    "interval",
    "path"
  ],
  "properties": {
    "interval": {
      "$id": "#/properties/interval",
      "type": "integer",
      "title": "The Interval Schema"
    },
    "path": {
      "$id": "#/properties/path",
      "type": "string",
      "title": "The Path Schema"
    },
    "repos": {
      "$id": "#/properties/repos",
      "type": "array",
      "title": "The Repos Schema",
      "items": {
        "$id": "#/properties/repos/items",
        "type": "object",
        "title": "The Items Schema",
        "required": [
          "url",
          "type"
        ],
        "properties": {
          "url": {
            "$id": "#/properties/repos/items/properties/url",
            "type": "string",
            "title": "The Url Schema"
          },
          "httpauth": {
            "$id": "#/properties/repos/items/properties/httpauth",
            "type": "object",
            "title": "The Httpauth Schema",
            "required": [
              "user",
              "token"
            ],
            "properties": {
              "user": {
                "$id": "#/properties/github/items/properties/httpauth/properties/user",
                "type": "string",
                "title": "The User Schema"
              },
              "token": {
                "$id": "#/properties/github/items/properties/httpauth/properties/token",
                "type": "string",
                "title": "The Token Schema"
              }
            }
          },
          "type": {
            "$id": "#/properties/repos/items/properties/type",
            "type": "string",
            "title": "The Type Schema",
            "enum": [
              "fetch_mirror",
              "push_mirror"
            ]
          },
          "path": {
            "$id": "#/properties/repos/items/properties/path",
            "type": "string",
            "title": "The Path Schema"
          },
          "refs": {
            "$id": "#/properties/repos/items/properties/refs",
            "type": [
              "array",
              "null"
            ],
            "title": "The Refs Schema",
            "items": {
              "$id": "#/properties/repos/items/properties/refs/items",
              "type": "string",
              "title": "The Items Schema"
            }
          },
          "metadata": {
            "$id": "#/properties/repos/items/properties/metadata",
            "type": [
              "array",
              "null"
            ],
            "title": "The Metadata Schema",
            "items": {
              "$id": "#/properties/repos/items/properties/metadata/items",
              "type": "string",
              "title": "The Items Schema",
              "default": "",
              "examples": [
                "cgit"
              ],
              "enum": [
                "cgit"
              ]
            }
          }
        }
      }
    },
    "github": {
      "$id": "#/properties/github",
      "type": "array",
      "title": "The Github Schema",
      "items": {
        "$id": "#/properties/github/items",
        "type": "object",
        "title": "The Items Schema",
        "required": [
          "username"
        ],
        "properties": {
          "username": {
            "$id": "#/properties/github/items/properties/username",
            "type": "string",
            "title": "The Username Schema"
          },
          "httpauth": {
            "$ref": "#/properties/repos/items/properties/httpauth"
          },
          "repos": {
            "$id": "#/properties/github/items/properties/repos",
            "type": "boolean",
            "title": "The Repos Schema"
          },
          "protocol": {
            "$id": "#/properties/github/items/properties/protocol",
            "type": "string",
            "title": "The Protocol Schema",
            "enum": [
              "local",
              "ssh",
              "http",
              "git",
              ""
            ]
          },
          "repotype": {
            "$id": "#/properties/github/items/properties/repotype",
            "type": "string",
            "title": "The Repotype Schema",
            "enum": [
              "all",
              "public",
              "private",
              "forks",
              "sources",
              ""
            ]
          },
          "metadata": {
            "$id": "#/properties/github/items/properties/metadata",
            "type": [
              "array",
              "null"
            ],
            "title": "The Metadata Schema",
            "items": {
              "$id": "#/properties/github/items/properties/metadata/items",
              "type": "string",
              "title": "The Items Schema",
              "enum": [
                "cgit"
              ]
            }
          }
        }
      }
    }
  }
}
`

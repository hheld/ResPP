module Routing exposing (..)

import Navigation
import UrlParser as Url exposing ((</>), top, s)


type Route
    = HomeRoute
    | OtherRoute


route : Url.Parser (Route -> a) a
route =
    Url.oneOf
        [ Url.map HomeRoute top
        , Url.map OtherRoute (s "otherRoute")
        ]

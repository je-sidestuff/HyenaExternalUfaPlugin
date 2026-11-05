module hyena-ufa-plugin

go 1.21

require representable v0.0.0

// representable is the wire contract with local-representative. It lives in
// AI-evo1-spiral; a copy is vendored here so this repo builds on its own.
replace representable => ./third_party/representable

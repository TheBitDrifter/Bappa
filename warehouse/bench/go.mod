module github.com/TheBitDrifter/bappa/warehouse/bench

go 1.24.1

require (
	github.com/TheBitDrifter/bappa/table v0.0.0
	github.com/TheBitDrifter/bappa/warehouse v0.0.0-00010101000000-000000000000
	github.com/mlange-42/arche v0.13.2
)

require (
	github.com/TheBitDrifter/bark v0.0.0-20250302175939-26104a815ed9 // indirect
	github.com/TheBitDrifter/mask v0.0.1-early-alpha.1 // indirect
	github.com/TheBitDrifter/util v0.0.0-20241102212109-342f4c0a810e // indirect
)

replace github.com/TheBitDrifter/bappa/warehouse => ../

replace github.com/TheBitDrifter/bappa/table => ../../table/

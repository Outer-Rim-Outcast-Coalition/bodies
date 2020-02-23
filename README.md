# bodies

When engineering our ships, most of us, will follow one the guides that point us
to the different [Crystalline Shards](https://elite-dangerous.fandom.com/wiki/Crystalline_Shard) sites.
From the Wiki:

> To make this more generally useful we would need to complete the list of materials and hopefully find locations closer to the Bubble.

That's what this program attempts to do. It uses data dumps from EDSM:

* bodies.json.gz
* systemsWithCoordinates.json.gz

## Filtering criteria

* Distance from arrival > 12,000 Ls
* Distance from Sol <= 2,000 ly
* Having some kind of volcanism (any kind will do)
* Being landable (obviously...)
* Having one of the following interesting materials on it:
  * Ruthenium
  * Antimony
  * Yttrium
  * Technetium
  * Polonium
  * Tellurium

## Usage

### Pre-compute distances

Run the `computeDistances` command first:

`bodies computeDistances -g distances.gob -s systemsWithCoordinates.json.gz`

This will create a file that contains a lookup table of EDSM system IDs to distances of the system from the reference point (currently Sol).
This file will be used by the next step to enrich the bodies data. This process takes around 5 minutes on a 2019 MacBook Pro.

### Filter the bodies

Next run:

`bodies filterBodies -b bodies.json.gz -g distances.gob -o candidates.csv -f csv -l 1000000`

This will import the lookup table created in the previous step and start processing the data dump.
It will result in a CSV file with the following fields:

* Body Name
* System Name
* Distance from Reference (Sol)
* Distance of body from arrival star
* Body surface gravity
* Body surface temperature

It will limit itself to processing 1,000,000 bodies from the input. Without specifying `-l` the whole file will be processed.
**BEWARE:** at the time of this writing, the dump contains more than 150 million bodies. On a 2019 MacBook Pro that is a >2hr process.

You may also opt for `-f json` for a JSON output instead. The JSON output will contain all the fields from the data dump 
plus a *Distance* field which is the system's distance from reference (Sol).

## Development

This was originally developed by CMDR sugoruyo of the [Outer-Rim Outcasts Coalition](https://inara.cz/squadron/2959/). I'm too lazy to actually be going out that far out of the bubble and flying 300kLs in SC.

ANY help with this is appreciated whether you wanna help refine the criteria, contribute bug reports and/or fixes, feature proposals etc.

## License

Released under GPL v3.0.
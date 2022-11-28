# pipelines

It is a bit opinionated:
* I don't want to mimic node's Promises.
* I want plain async logic that can run multiple IN/OUT channels and check at a glance that types are correct.
* Pipelines may send out-of-band notifications to a unique bus channel.

TODO: basic pipelines work, but it is still difficult to get a nice, idiomatic chain of pipelines running smooth as I had expected.

Pipelines patterns to support:
* in/out/bus
* fan-int
* fan-out
* 2-way hetereogenous join
* feeder (no input)
* collector (no output)

> Findings
> I realize the implications of the limitation that no method can be itself parametric:
> this totally prevents me from building a fluent pipeline chain with a method like `Then[NEWOUT](next *Pipeline[OUT,NEWOUT]) *ChainedPipeline[IN, NEWOUT]`
> Ugh. Need to reflect more on that one.

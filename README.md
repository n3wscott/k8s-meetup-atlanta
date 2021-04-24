# k8s-meetup-atlanta

## Rainbow Deploys demo

Getting to eventually a `spec.traffic` like:

```yaml
  traffic:
  - percent: 40
    revisionName: colors-00002
    tag: green
  - percent: 50
    revisionName: colors-00003
    tag: blue
  - percent: 10
    revisionName: colors-00001
    tag: red
```
---
title: Licenses
nav_order: 5
---
Not all licensing features have been rolled out yet, so let me briefly explain the overall vision.

DevTUI will be available to everyone as a free product, with no limitations on functionality. However, it will include a nagging popup prompting users to support further development. Users can remove this nag screen by paying a (reasonable) one-time fee.

One license is currently valid for use with 3 machines.

[Customer portal](https://polar.sh/krooni/portal/request){: .btn .fs-5 .mb-4 .mb-md-0 .mr-2 } [Buy license](https://buy.polar.sh/polar_cl_JPBTnQKWsNBC8lA7tpR1uZYne5hMuW40xqTRI3P9WcH){: .btn .fs-5 .mb-4 .mb-md-0 .mr-2}
{: .text-center}
## Activation
It's possible to activate a license only through application.

`devtui license activate --key=DEVTUI-2CA57A34-E191-4290-A394-XXXXXX`

## Validation
It's possible to check that license was validated through application or through customer portal.

`devtui license validate`

## Deactivation
License deactivation happens through application or customer portal.

`devtui license deactivate`

Using application for deactivating might be useful for two cases:
- Running the app in CI environments where machines constantly get recycled
- Switching from one machine to another

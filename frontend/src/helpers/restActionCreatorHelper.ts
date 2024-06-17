export function restActionCreatorHelper<T extends string>(reducer: T) {
    return function <D extends string>(restAction: D) {
       return {
          request: `${reducer}/${restAction}_request`,
          success: `${reducer}/${restAction}_success`,
          failure: `${reducer}/${restAction}_failure`,
          needUpdate: `${reducer}/${restAction}_needUpdate`,
       } as const
    }
 }
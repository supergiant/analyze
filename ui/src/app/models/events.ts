// TODO: types in this file need to be shared between Analyze Core and all plugins
export enum EventType {
  CE_LOADED_EVENT = "CELoadedEvent",
  ACTION_SUBMIT_EVENT = "ActionSubmitEvent",
}

export interface CELoadedEvent {
  pluginId: string
  pluginVersion: string
  webComponentName: string
  selector: string
}

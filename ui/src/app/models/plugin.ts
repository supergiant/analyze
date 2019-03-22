export interface Plugin {
  checkComponentEntryPoint: string;
  description: string;
  id: string;
  installedAt: string;
  name: string;
  serviceEndpoint: string;
  serviceLabels: any,
  settingsComponentEntryPoint: string;
  status: string;
  version: string;
}

import React, { Component, ReactElement } from 'react';

import {
  BurndownType,
  RepoBurndown,
  GlobalBurndown,
} from 'lib';
import { Error } from 'components/alerts';

import { globalBurndownElement } from './global';

interface BurndownProps {
  data: RepoBurndown;
}

class Burndown extends Component<BurndownProps> {
  public render(): ReactElement {
    const { data } = this.props;
    switch (data.type) {
      case BurndownType.GLOBAL: return globalBurndownElement(data as GlobalBurndown);
      case BurndownType.FILE: return <Error message="unimplemented" />;
      case BurndownType.AUTHOR: return <Error message="unimplemented" />;
      case BurndownType.ALERT: return <Error message={`error: ${data}`} />;
      default: return <Error message="invalid data found" />;
    }
  }
}

export default Burndown;

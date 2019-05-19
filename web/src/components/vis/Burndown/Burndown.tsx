import React, { Component, ReactElement } from 'react';
import DeckGL, { ScatterplotLayer } from 'deck.gl';

import { Error } from '../../alerts';

import { Repo_repo_burndown, Repo_repo_burndown_GlobalBurndown } from '../../../lib/queries/types/Repo'; // eslint-disable-line
import { BurndownType } from '../../../lib/queries/types/global';

// Alias types for readability
type RepoBurndown = Repo_repo_burndown; // eslint-disable-line
type GlobalBurndown = Repo_repo_burndown_GlobalBurndown; // eslint-disable-line

// Viewport settings
const viewState = {
  longitude: -74,
  latitude: 40.76,
  zoom: 13,
  maxZoom: 16,
  pitch: 50,
  bearing: 50,
};

const MALE_COLOR = [0, 128, 255];
const FEMALE_COLOR = [255, 0, 128];

function globalBurndownElement(data: GlobalBurndown): ReactElement {
  const layers = [
    new ScatterplotLayer({
      id: 'scatter-plot',
      data: 'https://raw.githubusercontent.com/uber-common/deck.gl-data/master/examples/scatterplot/manhattan.json',
      radiusScale: 10,
      radiusMinPixels: 0.5,
      getPosition: (d: number[]): number[] => [d[0], d[1], 0],
      getColor: (d: number[]): number[] => (d[2] === 1 ? MALE_COLOR : FEMALE_COLOR),
    }),
  ];

  return (
    <DeckGL viewState={viewState} layers={layers} />
  );
}

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

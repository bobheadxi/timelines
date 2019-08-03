import React, { ReactElement, Component } from 'react';
import DeckGL, { HexagonLayer } from 'deck.gl';

import { GlobalBurndown, GlobalBurndownEntry } from 'lib';
import { Coordinate } from 'components/vis';

class Global extends Component<{
  data: GlobalBurndown;
}> {
  // https://deck.gl/#/documentation/developer-guide/adding-interactivity?section=using-react
  private renderTooltip() {
    const { hoveredObject, pointerX, pointerY } = this.state || {};
    return hoveredObject && (
      <div style={{position: 'absolute', zIndex: 1, pointerEvents: 'none', left: pointerX, top: pointerY}}>
        { hoveredObject.message }
      </div>
    );
  }

  // https://deck.gl/#/documentation/deckgl-api-reference/layers/hexagon-layer
  public render(): ReactElement {
    const { viewState } = this.state;

    const layers = [
      new HexagonLayer({
        id: 'scatter-plot',
        data: 'https://raw.githubusercontent.com/uber-common/deck.gl-data/master/website/sf-bike-parking.json',
        extruded: true,
        radius: 200,
        elevationScale: 4,
        getPosition: d => d.COORDINATES,
        onHover: ({ object, x, y }): void => {
          const tooltip = `${object.centroid.join(', ')}\nCount: ${object.points.length}`;
          /* Update tooltip
             http://deck.gl/#/documentation/developer-guide/adding-interactivity?section=example-display-a-tooltip-for-hovered-object
          */
        },
      }),
    ];

    return (
      <DeckGL viewState={viewState} layers={layers} />
    );
  }
}

function coordinates(entry: GlobalBurndownEntry): Coordinate {
  return [0, 0, 0];
}

function elevation(): number {
  return 0;
}

export default Global;

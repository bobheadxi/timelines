import React, { Component, ReactElement } from 'react';

export interface CardButton {
  href: string;
  text: string;
}

export interface Card {
  title: string;
  body: string;
  button?: CardButton;
}

class CardSet extends Component<{
  cards: Card[];
}> {
  public render(): ReactElement {
    const { cards } = this.props;
    return (
      <div
        className="uk-child-width-1-2@s uk-grid-match"
        data-uk-scrollspy="target: > div; cls:uk-animation-fade; delay: 50"
        data-uk-grid
      >
        {cards.map((card): ReactElement => {
          const { title, body, button } = card;
          return (
            <div>
              <div className="uk-card uk-card-hover uk-card-default">
                <div className="uk-card-body">
                  <h3 className="uk-card-title">
                    {title}
                  </h3>
                  <p>
                    {body}
                  </p>
                  {button
                    ? (
                      <a href={button.href} className="uk-button uk-button-text">
                        {button.text}
                      </a>
                    ) : null}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    );
  }
}

export default CardSet;

module.exports = {
  parser: '@typescript-eslint/parser',
	plugins: ['@typescript-eslint', 'jest'],
	extends: [
		'react-app',
		'airbnb',
		'plugin:@typescript-eslint/recommended',
		'plugin:jest/recommended',
	],
	env: {
		browser: true
	},
	rules: {
		'react/jsx-filename-extension': [1, { 'extensions': ['.tsx'] }],
		'@typescript-eslint/indent': ['error', 2],

		// TODO: re-enable these when project is more stable
		'react/prefer-stateless-function': false,
		'import/prefer-default-export': false,

		// TODO: not really sure what these are, related to forms
		'jsx-a11y/label-has-for': false,
		'jsx-a11y/label-has-associated-control': false,
	},
	settings: {
		'import/resolver': {
			node: {
				paths: ['src'],
				extensions: [
					'.js',
					'.jsx',
					'.ts',
					'.tsx',
				],
			},
		},
	},
}

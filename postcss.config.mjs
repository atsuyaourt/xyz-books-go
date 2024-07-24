import UnoCSS from '@unocss/postcss'
import PostCSSImport from 'postcss-import'

export default {
  plugins: [UnoCSS(), PostCSSImport()],
}

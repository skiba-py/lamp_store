import { Grid } from '@chakra-ui/react';
import ProductCard from './ProductCard';

export default function ProductList({ products, onAddToCart }) {
  return (
    <Grid templateColumns="repeat(auto-fit, minmax(220px, 1fr))" gap={8}>
      {products.map(product => (
        <ProductCard key={product.id} product={product} onAddToCart={() => onAddToCart(product)} />
      ))}
    </Grid>
  );
}

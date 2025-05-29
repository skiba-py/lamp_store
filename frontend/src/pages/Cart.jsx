import { Box, Heading, Text, Flex, Spinner, useToast } from '@chakra-ui/react';
import CartList from '../components/CartList';
import { useEffect, useState } from 'react';
import { ordersApi } from '../services/api';

const USER_ID = "123e4567-e89b-12d3-a456-426614174000";

export default function Cart() {
  const [order, setOrder] = useState(null);
  const [loading, setLoading] = useState(true);
  const [updating, setUpdating] = useState(false);
  const toast = useToast();

  useEffect(() => {
    const fetchCart = async () => {
      setLoading(true);
      try {
        const cart = await ordersApi.getOrders(USER_ID);
        setOrder(cart || null);
      } catch (err) {
        toast({ title: 'Ошибка', description: 'Не удалось загрузить корзину', status: 'error' });
      } finally {
        setLoading(false);
      }
    };
    fetchCart();
  }, []);

  const handleUpdateQuantity = async (itemId, delta) => {
    if (!order) return;
    setUpdating(true);
    try {
      const updatedItems = order.items.map(item =>
        item.id === itemId ? { ...item, quantity: Math.max(1, item.quantity + delta) } : item
      );
      const itemsForBackend = updatedItems.map(item => ({
        id: item.id,
        order_id: item.order_id,
        product_id: item.product.id,
        quantity: item.quantity,
        price: item.product.price
      }));
      const updatedOrder = { ...order, items: itemsForBackend };
      await ordersApi.updateOrder(order.id, updatedOrder);
      setOrder({ ...order, items: updatedItems });
    } catch (err) {
      toast({ title: 'Ошибка', description: 'Не удалось обновить корзину', status: 'error' });
    } finally {
      setUpdating(false);
    }
  };

  if (loading) return <Box textAlign="center" py={10}><Spinner size="xl" /></Box>;

  if (!order || !order.items || order.items.length === 0) {
    return (
      <Box id="cart" w="100%" minH="100vh" maxW="600px">
        <Heading as="h1" size="lg" mb={6}>Корзина</Heading>
        <Text color="gray.500">Корзина пуста</Text>
      </Box>
    );
  }

  const total = order.items.reduce((sum, item) => sum + (item.product.price || 0) * item.quantity, 0);

  return (
    <Box id="cart" w="100%" minH="100vh" maxW="800px">
      <Heading as="h1" size="lg" mb={6}>Корзина</Heading>
      <CartList
        items={order.items}
        onRemove={async (id) => {
          if (!order) return;
          setUpdating(true);
          try {
            const updatedItems = order.items.filter(item => item.id !== id);
            if (updatedItems.length === 0) {
              await ordersApi.deleteOrder(order.id);
              setOrder(null);
              toast({ title: 'Корзина очищена', description: 'Все товары удалены', status: 'info' });
            } else {
              const itemsForBackend = updatedItems.map(item => ({
                id: item.id,
                order_id: item.order_id,
                product_id: item.product.id,
                quantity: item.quantity,
                price: item.product.price
              }));
              const updatedOrder = { ...order, items: itemsForBackend };
              await ordersApi.updateOrder(order.id, updatedOrder);
              setOrder({ ...order, items: updatedItems });
            }
          } catch (err) {
            toast({ title: 'Ошибка', description: 'Не удалось удалить товар из корзины', status: 'error' });
          } finally {
            setUpdating(false);
          }
        }}
        onUpdateQuantity={handleUpdateQuantity}
        updating={updating}
      />
      <Box mt={6} textAlign="right">
        <Text fontWeight="bold" fontSize="xl">Итого: {total.toFixed(2)} ₽</Text>
      </Box>
    </Box>
  );
}

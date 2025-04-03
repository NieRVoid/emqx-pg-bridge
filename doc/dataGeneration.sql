-- Insert 6 rooms
INSERT INTO rooms (number, name, description, occupancy)
VALUES
('101', 'Room 101', 'First floor room 101', 'unknown'),
('102', 'Room 102', 'First floor room 102', 'unknown'),
('201', 'Room 201', 'Second floor room 201', 'unknown'),
('202', 'Room 202', 'Second floor room 202', 'unknown'),
('301', 'Room 301', 'Third floor room 301', 'unknown'),
('302', 'Room 302', 'Third floor room 302', 'unknown');

-- Insert the room_status corresponding to 6 rooms, all rooms are empty rooms
INSERT INTO room_status (room_id, temperature, humidity, air_quality, light_level, noise_level, occupied, occupant_count, count_confidence, occupied_confidence, count_source, last_source_change, metadata)
SELECT id, 22, 40, 80, 300, 30, false, 0, 0, 0, 'system', now(), '{}'::jsonb
FROM rooms;

-- Insert 6 devices corresponding to each room
DO $$
DECLARE
  r RECORD;
  i INTEGER;
  device_uuid UUID;
  device_id INTEGER;
  led_status JSONB := '{"power": "off", "color": {"r": 0, "g": 0, "b": 0}}';
BEGIN
  FOR r IN SELECT id, number FROM rooms LOOP
    FOR i IN 1..6 LOOP
      device_uuid := gen_random_uuid();
      INSERT INTO devices (uuid, name, type, model, manufacturer, description, room_id)
      VALUES (
        device_uuid,
        'LED-' || r.number || '-' || i,
        'RGB-LED',
        'Model X',
        'LED Manufacturer',
        'RGB LED device',
        r.id
      ) RETURNING id INTO device_id;

      -- Insert the initial closed state for each device
      INSERT INTO device_status (device_id, status, updated_at, last_reported_at)
      VALUES (
        device_id,
        led_status,
        now(),
        now()
      );
    END LOOP;
  END LOOP;
END;
$$;

/**
 * Builds the photo-matched stylized hero as a plain THREE.Group (no React).
 * Proportions tuned for a full-body hero frame (head must stay inside the canvas).
 */
import * as THREE from 'three'

function mat(color: string, roughness: number, metalness = 0) {
  return new THREE.MeshStandardMaterial({
    color,
    roughness,
    metalness,
    envMapIntensity: 0.35,
  })
}

export function buildHeroCharacterGroup(): THREE.Group {
  const root = new THREE.Group()
  root.name = 'HeroRoot'

  const materials = {
    skin: mat('#c99472', 0.62, 0.02),
    hair: mat('#0a0a0b', 0.88),
    brow: mat('#111113', 0.92),
    shirt: mat('#9aa1aa', 0.82),
    pants: mat('#6b82a0', 0.75),
    shoe: mat('#e8e2d8', 0.58),
    shoeAccent: mat('#2dd4bf', 0.48),
    sole: mat('#f1f5f9', 0.7),
    eye: mat('#1a120e', 0.4),
    smile: mat('#b45c4c', 0.55),
  }

  const add = (
    parent: THREE.Object3D,
    geo: THREE.BufferGeometry,
    material: THREE.Material,
    opts?: {
      name?: string
      position?: [number, number, number]
      rotation?: [number, number, number]
      scale?: [number, number, number]
    },
  ) => {
    const mesh = new THREE.Mesh(geo, material)
    mesh.castShadow = true
    mesh.receiveShadow = true
    if (opts?.name) mesh.name = opts.name
    if (opts?.position) mesh.position.set(...opts.position)
    if (opts?.rotation) mesh.rotation.set(...opts.rotation)
    if (opts?.scale) mesh.scale.set(...opts.scale)
    parent.add(mesh)
    return mesh
  }

  // Legs — slim, light wash jeans
  add(root, new THREE.CapsuleGeometry(0.1, 0.55, 6, 14), materials.pants, {
    name: 'LegL',
    position: [-0.17, 0.42, 0],
  })
  add(root, new THREE.CapsuleGeometry(0.1, 0.55, 6, 14), materials.pants, {
    name: 'LegR',
    position: [0.17, 0.42, 0],
  })

  for (const x of [-0.17, 0.17] as const) {
    const shoe = new THREE.Group()
    shoe.name = x < 0 ? 'ShoeL' : 'ShoeR'
    shoe.position.set(x, 0.05, 0.05)
    add(shoe, new THREE.BoxGeometry(0.2, 0.09, 0.32), materials.shoe, {
      rotation: [0.12, 0, 0],
    })
    add(shoe, new THREE.BoxGeometry(0.22, 0.045, 0.34), materials.sole, {
      position: [0, -0.06, 0.01],
    })
    add(shoe, new THREE.BoxGeometry(0.025, 0.035, 0.14), materials.shoeAccent, {
      position: [0.1, 0.01, 0],
    })
    root.add(shoe)
  }

  // Torso — rounded box silhouette (less “balloon capsule”)
  const torso = new THREE.Group()
  torso.name = 'Torso'
  torso.position.set(0, 1.02, 0)

  add(torso, new THREE.CapsuleGeometry(0.28, 0.42, 8, 18), materials.shirt, {
    name: 'Shirt',
    scale: [1.05, 1, 0.92],
  })
  // Slight oversized hem
  add(torso, new THREE.CylinderGeometry(0.3, 0.32, 0.14, 18), materials.shirt, {
    position: [0, -0.38, 0],
  })
  add(torso, new THREE.CapsuleGeometry(0.11, 0.14, 6, 12), materials.shirt, {
    position: [-0.34, 0.22, 0],
    rotation: [0, 0, 0.65],
  })
  add(torso, new THREE.CapsuleGeometry(0.11, 0.14, 6, 12), materials.shirt, {
    position: [0.34, 0.22, 0],
    rotation: [0, 0, -0.65],
  })

  // Neck
  add(torso, new THREE.CylinderGeometry(0.08, 0.09, 0.12, 14), materials.skin, {
    name: 'Neck',
    position: [0, 0.42, 0],
  })

  const leftArm = new THREE.Group()
  leftArm.name = 'LeftArm'
  leftArm.position.set(-0.42, 0.1, 0)
  add(leftArm, new THREE.CapsuleGeometry(0.075, 0.38, 6, 12), materials.skin, {
    position: [0, -0.28, 0],
  })
  torso.add(leftArm)

  const rightArm = new THREE.Group()
  rightArm.name = 'RightArm'
  rightArm.position.set(0.42, 0.1, 0)
  add(rightArm, new THREE.CapsuleGeometry(0.075, 0.38, 6, 12), materials.skin, {
    position: [0, -0.28, 0],
  })
  torso.add(rightArm)

  const head = new THREE.Group()
  head.name = 'Head'
  head.position.set(0, 0.62, 0)

  add(head, new THREE.SphereGeometry(0.24, 36, 36), materials.skin, {
    name: 'Face',
    scale: [0.9, 1.02, 0.92],
  })

  // Hair volume (black, messy top + fringe)
  add(head, new THREE.SphereGeometry(0.25, 28, 28), materials.hair, {
    position: [0, 0.14, -0.03],
    scale: [1.12, 0.72, 1.15],
  })
  add(head, new THREE.SphereGeometry(0.2, 24, 24), materials.hair, {
    position: [0, 0.26, 0],
    scale: [1, 0.55, 1],
  })
  add(head, new THREE.SphereGeometry(0.14, 16, 16), materials.hair, {
    position: [-0.07, 0.12, 0.16],
    rotation: [0.45, 0.15, -0.2],
    scale: [0.7, 0.4, 0.55],
  })
  add(head, new THREE.SphereGeometry(0.14, 16, 16), materials.hair, {
    position: [0.08, 0.13, 0.15],
    rotation: [0.4, -0.2, 0.25],
    scale: [0.65, 0.38, 0.5],
  })
  add(head, new THREE.SphereGeometry(0.12, 16, 16), materials.hair, {
    position: [0, 0.08, 0.18],
    rotation: [0.7, 0, 0],
    scale: [0.85, 0.35, 0.45],
  })

  // Brows + eyes + upward smile
  add(head, new THREE.BoxGeometry(0.085, 0.022, 0.025), materials.brow, {
    position: [-0.08, 0.07, 0.21],
    rotation: [0.05, 0, 0.18],
  })
  add(head, new THREE.BoxGeometry(0.085, 0.022, 0.025), materials.brow, {
    position: [0.08, 0.07, 0.21],
    rotation: [0.05, 0, -0.18],
  })
  add(head, new THREE.SphereGeometry(0.028, 14, 14), materials.eye, {
    position: [-0.07, 0.02, 0.215],
  })
  add(head, new THREE.SphereGeometry(0.028, 14, 14), materials.eye, {
    position: [0.07, 0.02, 0.215],
  })
  // Smile opens upward (rotate so arc faces up)
  add(head, new THREE.TorusGeometry(0.055, 0.01, 10, 24, Math.PI), materials.smile, {
    position: [0, -0.055, 0.21],
    rotation: [0, 0, Math.PI],
    scale: [1, 0.55, 0.55],
  })

  torso.add(head)
  root.add(torso)

  // Origin at feet so camera framing is predictable
  root.position.set(0, 0, 0)
  root.scale.setScalar(1)
  return root
}
